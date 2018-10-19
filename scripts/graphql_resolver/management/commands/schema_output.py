import sys, inspect
import io, yaml
import argparse
import stringcase

from collections import OrderedDict

from django.db import models
from django.core.management.base import BaseCommand, CommandError

from candidate.models import *

class Command(BaseCommand):
    help = 'Generates meta about model classes'

    def add_arguments(self, parser):
        parser.add_argument(
            'outfile',
            nargs='?',
            type=argparse.FileType('w'),
            default=sys.stdout,
            help='output file for generated metadata')

    def handle(self, *args, **options):
        outfile = options["outfile"]

        m = self._get_model_metadata()
        # print json.dumps(m, indent=2)
        # outfile.write(json.dumps(m, indent=2, ensure_ascii=False))
        outfile.write(yaml.dump({"models": m}, default_flow_style=False))
        outfile.write("\n")

    def _get_model_metadata(self):
        m = []
        for name, obj in self._get_model_classes():
            fields = []

            unique = self._get_unique_together(obj)

            print "Model %s" % name
            for f in obj._meta.get_fields():
                if f.is_relation:
                    if f.many_to_many:
                        through = getattr(f.remote_field, "through", "")
                        if not through:
                            continue

                        indexes = []
                        for u in unique:
                            if str(f.name) in u:
                                indexes.append(u)

                        fields.append({
                            "name": str(f.name),
                            "relationship": {
                                "through": stringcase.snakecase(through._meta.object_name),

                                "to": stringcase.snakecase(f.related_model._meta.object_name),
                                "field": stringcase.snakecase(name + "Id"),
                                "type": "many2many",
                            },
                            "indexes": indexes,
                        })

                    elif f.many_to_one:
                        fields.append({
                            "name": str(f.name),
                            "internal": str(f.column),
                            "type": "id",
                        })

                    elif f.one_to_one:
                        indexes = []
                        for u in unique:
                            if str(f.name) in u:
                                indexes.append(u)

                        fields.append({
                            "name": str(f.name),
                            "internal": str(f.column),
                            "relationship": {
                                "to": stringcase.snakecase(f.related_model._meta.object_name),
                                "type": "one2one",
                            },
                            "indexes": indexes,
                        })

                    elif f.one_to_many:
                        fields.append({
                            "name": str(f.name),
                            "relationship": {
                                "to": stringcase.snakecase(f.related_model._meta.object_name),
                                "field": stringcase.snakecase(name + "Id"),
                                "type": "one2many",
                            },
                        })

                else:
                    indexes = []
                    if f.unique:
                        if str(f.name) == "id":
                            indexes.append("primary")
                        else:
                            indexes.append(str(f.name + "_unique"))
                    for u in unique:
                        if str(f.name) in u:
                            indexes.append(u)

                    fields.append({
                        "name": str(f.name),
                        "internal": str(f.column),
                        "type": self._get_type(f),
                        "indexes": indexes,
                    })

            m.append({
                "name": stringcase.snakecase(name),
                "description": name,
                "internal": self._get_model_table(obj),
                "fields": fields,
                "mutations": self._get_model_mutations(obj),
            })

        return m

    def _get_model_table(self, obj):
        return str(obj._meta.db_table)

    def _get_unique_together(self, obj):
        u = []
        for tuple in obj._meta.unique_together:
            u.append(str("_".join(tuple) + "_unique"))
        return u

    def _get_model_mutations(self, obj):
        fields = []
        for f in obj._meta.get_fields(include_hidden=False):
            if not f.many_to_many and not f.one_to_many and not f.one_to_one:
                fields.append(str(f.name))

        arr = []
        arr.append({
            "name": "create",
            "type": "insert",
            "fields": filter(lambda x: x != "id", fields),
        })
        arr.append({
            "name": "update",
            "type": "update",
            "fields": fields,
        })
        arr.append({
            "name": "delete",
            "type": "delete",
            "fields": ["id"],
        })
        return arr

    def _get_model_classes(self):
        return inspect.getmembers(sys.modules[__name__], ismodel)

    def _get_type(self, f):
        lookup = {
            "DateTimeField": "timestamp",
            "AutoField": "id",
            "CharField": "string",
            "BooleanField": "boolean",
            "TextField": "string",
            "IntegerField": "integer",
            "FloatField": "float",
            "EmailField": "email",
            "ForeignKey": "id", # "one2many",
            "OneToOneField": "one2one",
            "ManyToManyField": "many2many",
        }
        t = f.get_internal_type()
        return lookup.get(t, "unknown")

def isattribute(obj):
    return isinstance(obj, models.query_utils.DeferredAttribute)

def ismodel(obj):
    if inspect.isclass(obj):
        if isinstance(obj, models.base.ModelBase):
            return obj
