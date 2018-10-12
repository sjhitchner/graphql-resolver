import sys, inspect
import io, json
import argparse
from collections import OrderedDict

from django.db import models
from django.core.management.base import BaseCommand, CommandError

from slack.models import *


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

        # self._get_model_metadata()
        m = self._get_model_metadata()
        # print json.dumps(m, indent=2)
        # outfile.write(json.dumps(m, indent=2, ensure_ascii=False))
        # outfile.write(yaml.dump(self._get_model_metadata(), default_flow_style=False))
        outfile.write("\n")

    def _get_model_metadata(self):
        m = []
        for name, cls in self._get_model_classes():
            m.append({
                "name": name,
                "description": "",
                "internal": self._get_model_table(cls),
                # "fields": self._get_model_fields(name, cls),
                "mutations": self._get_model_mutations(cls),
            })

            obj = cls
            print "====="
            print "Model %s" % name
            # print dir(obj)
            #print dir(obj._meta)
            print obj._meta.fields
            for f in obj._meta.fields:
                print f.name, f.column, f.get_internal_type()

                if f.is_relation:
                    print f.related_model
                    print f.foreign_related_fields
                    if f.one_to_one:
                        print "O2O"
                    if f.many_to_one:
                        print "M2O"
                    if f.many_to_many:
                        print "M2M"
                print dir(f)
            print "====="

        return m

    def _get_model_table(self, obj):
        return str(obj._meta.db_table)

    def _get_model_fields(self, name, obj):
        print "Model %s" % name

        print "====="
        print dir(obj)
        print dir(obj._meta)
        print obj._meta.fields
        for f in obj._meta.fields:
            print f.name, f.column, f.get_internal_type(),

            print dir(f)
        print "====="

        arr = []
        for f in obj._meta.get_fields(include_hidden=True):
            indexes = []
            if f.name == "id":
                indexes.append("primary")

            field = {
                "name": str(f.name),
                "description": "",
                "expose": True,
                "indexes": indexes,
            }

            t = self._get_type(f)
            if t == "one2one":
                field["relationship"] = {
                    "to": f.related_model._meta.model_name,
                    "type": t,
                }
            elif t == "one2many":
                field["relationship"] = {
                    "to": f.related_model._meta.model_name,
                    "field": name.lower() + "_id",
                    "type": t,
                }
            elif t == "many2many":
                print "A", dir(f)
                print "A", f.name
                print "B", f.get_related_field()
                print "B", dir(f.get_related_field())
                print f.rel.get_related_field()
                field["relationship"] = {
                    "to": f.related_model._meta.model_name,
                    "through": f.rel.through._meta.db_table,
                    "field": name.lower() + "_id",
                    "type": t,
                }
            else:
                field["type"] = t

            arr.append(field)
        return arr

    def _get_model_mutations(self, obj):
        fields = []
        for f in obj._meta.get_fields(include_hidden=False):
            if not f.many_to_many and not f.one_to_many and not f.one_to_one:
                fields.append(str(f.column))

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
            "ForeignKey": "one2many",
            "DateTimeField": "timestamp",
            "AutoField": "id",
            "CharField": "string",
            "OneToOneField": "one2one",
            "BooleanField": "boolean",
            "TextField": "string",
            "IntegerField": "integer",
            "FloatField": "float",
            "EmailField": "email",
            "ManyToManyField": "many2many",
        }

        t = f.get_internal_type()
        print "++>", t, f.name, f.get_internal_type()
        return lookup.get(t, "unknown")

def isattribute(obj):
    return isinstance(obj, models.query_utils.DeferredAttribute)

def ismodel(obj):
    if inspect.isclass(obj):
        if isinstance(obj, models.base.ModelBase):
            return obj
