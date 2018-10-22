# graphql-resolver

Utilities to generate the basic code necessary to create GraphQL Go service.

Example command:

	./graphql-resolver --config <path to schema files> --path <generation path>


## Configuration

To start a new service you need to define a schema file that define all models, 
relationships, custom Go types and mutations your service needs to support

An [example schema file](https://github.com/sjhitchner/graphql-resolver/blob/master/example/models.yml) is included in the examples directory of this project.  Running the below command with generate the code necessary to start the service

	./graphql-resolver --config example/models.yml --path ./example

## Schema File

The schem file is a yaml document. It can be divided into several files.  It is recommended you create a config subdirectory to hold all config files.  When the generator is run it looks in the config directory and concatanates all *.yml files into a single document.

### General Options
General options are a work in progress.  They are intended to configure global settings

	generate: 
	  - sql
	  - graphql
	  - resolvers

	graphql:
	  package: "resolvers"

	sql:
	  package: "db"
	  dialect: "postgres"

	resolvers:
	  package: "resolvers"


### Custom Types

Custom types all you to generate custom Go types that you can attach specific validation.  Custom types also
make your internal APIs clean as it allows you to create specific types like Username, Password instead of
just using strings.  If you add a `Validate() error` method to each type the Validate method will
automatically be run to validate values passed in are valid e.g. Valid username, valid email etc.

	custom_types:
	  - name: id
	    primative: integer
	  - name: username
	    primative: string
	  - name: email
	    primative: string
	  - name: password
	    primative: string


### Models

Models are where you define the objects used in your service.  The definition of each models will be generated
in the domain directory.

All names and field names should be in snake case (lowercase and underscores) e.g. my_name, my_field


#### Model

Models contain a name, description, list of fields and supported mutations.  

	models:
	  - name: <model name>
	    description: <model description>
	    fields: <list of fields>
	    mutations: <list of mutations>

#### Field

Fields contain name, optional internal name, description, expose, deprecated, type, relationship and indexes

* expose: whether the field is exposed publically
* deprecated: optional flag to indicate whether field is deprecated, strictly informational
* type: field type valid options
* * id
* * integer
* * float
* * string
* * custom type definted above
* * no type if relationship
* indexes: How this field is exposed in SQL.  To define an index use either `primary` or the field name and append either `_unique` or `_index`, for multi field append fieldname together with underscore.  
* * <field_name>_unique: defines a unique index where only one possible value is returned when queried
* * <field_name>_index: defines an index where multiple possible values are returned when queried
* * primary: defines a primary index
* relationship

	
##### Relationship

Relationships define how models work together

* to: model the relationship references, to field is always id
* field: field in the to model only needed for one2many and many2many
* through: the linking table to join through, only needed for many2many
* type: relationship type either one2one, one2many, many2many

###### One To One Relationship

Defines a one to one relations

	relationship:
      to: <to model>
      type: one2one
 
###### One To Many Relationship

    relationship:
      to: <to model>
      field: <to field>
      type: one2many

###### Many To Many Relationship

    relationship:
      to: <to model>
      field: <to field>
      through: <intermediary/linking table>
      type: "many2many"


##### Example

    fields:
      - name: "id"
	    description: "The ID of the User"
	    expose: true
	    deprecated: false
	    type: id 
	    indexes: 
	      - primary
	
      - name: "username"
	    description: "The Username of the User"
	    expose: true
	    type: username
	    indexes:
	      - username_unique
	
	  - name: "email"
	    description: "Email of the User"
	    expose: true
	    type: email
	    indexes:
	      - email_unique
	
	  - name: "password"
	    description: "Password of the User"
	    type: password

      - name: "team_list"
        description: "Teams users"
        expose: true
        relationship:
          to: "team"
          through: "team_member"
          field: "user_id"
          type: "many2many"

##### Mutation

Mutations are how the model is modified.  There are three valid fields, name, type, fields.
The name can be any name you choose in snake case. Three types of mutations are supported:
insert, update, delete.  The field list are the fields that will be modified by the mutation.

    mutations:
	  - name: create
	    type: insert
	    fields:
	      - username
	      - email
	      - password
	  - name: update
	    type: update
	    fields: 
	      - id
	      - username
	      - email
	      - password
	    key: id
	  - name: delete
	    type: delete
	    fields:
	      - id
	
# Customization

All files are generated with a `.gen.go` extension.  Currently, there is not logic in place generate 
changes when the generate is run again after adding new models.  If you need to specialize a file it is 
recommended that you rename the file to remove the `.gen.go` extension so a if the generator is run in
future you can simply delete the generated file.


# Motivation

When developing there is a large amount of boilerplate code needed to get everything working.  This library attempts
to automatically generate the basic framework necessary to get a working service and making it easy to extend as you 
see fit.  

The service is generated using a [Go Clean Code](https://manuel.kiessling.net/2012/09/28/applying-the-clean-architecture-to-go-applications/) structure. It uses the [Graph Gophers](https://github.com/graph-gophers) graphql library which is easy to use and seems to be fully featured.  

Some of the components were constructed using architecture ideas taken from [Oscar Yuen Go Graphql Starter Example](http://github.com/OscarYuen/go-graphql-starter).

# Features

Not all Graphql constructs are supported.  The schema definition does not support adding custom GraphQL types nor does it 
support ENUMs or Object Inheritance.




# Future Work

* Generate Unit Tests for 100% code coverage
* Allow integration for other backends (only SQL based/Postgres is currently supported)
* Support additional graphql constructs
** ENUMS
** Object Inheritance
* Additional tools to automatically generate schema file from additional sources
* Logic to only generate delta/changes
* More thorough validation





