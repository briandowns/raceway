# Raceway

Raceway is an application to manage [OpenStack](https://www.openstack.org/) [Rally](http://rally.readthedocs.org/en/latest/) benchmark runs.  Interaction with Raceway is driven through a REST API which feeds a frontend GUI and a CLI (whenever that gets written).

## Requirements

* OpenStack Rally
* MySQL

OpenStack Rally should be installed with the following steps.  This makes sure you are running the latest version.

```bash
$ wget https://github.com/openstack/rally/archive/master.zip
$ unzip master.zip
$ cd rally-master
$ ./install_rally.sh --system ../rally --overwrite --verbose --dbtype mysql --db-host localhost --db-user rally --db-password rally --db-name rally
```

If you have a different version of Python installed, use the example below and replace 'python2.7' with your version of Python.

```bash
$ ./install_rally.sh --target ../rally --overwrite --verbose --dbtype mysql --db-host localhost --db-user rally --db-password rally --db-name rally --python `which python2.7`
```

## Installation 

Simply put the Raceway binary in your path with the public directory.  Running `make build` gets all dependencies and builds the binary while also checking for race conditions.

## Usage

```bash
$ git clone git@github.com:briandowns/raceway.git
$ cd raceway
$ make build
$ ./raceway --help
```

### Web UI

An incomplete web UI is provided to be able to view the generated HTML or JSON from the task runs.

### REST API

| Method | Resource | Description | 
| :----- | :------- | :---------- |
| GET    | /api/v1/deployments                    | Retrieve a list of Rally deployments |
| GET    | /api/v1/deployments/{name}             | Retrieve the details of a given deployment |
| POST   | /api/v1/deployments                    | Create a new deployment |
| GET    | /api/v1/scenarios                      | Retrieve a list of scenarios |
| GET    | /api/v1/scenarios/{name}               | Retrieve the details of a given scenario |
| POST   | /api/v1/scenarios                      | Create a new scenario | 
| GET    | /api/v1/tasks                          | Retrieve a list of tasks |
| POST   | /api/v1/tasks                          | Create a new task |
| GET    | /api/v1/tasks/start                    | Start a given task |
| GET    | /api/v1/tasks/results/{task_uuid}      | Retrieve the results of a given task |
| GET    | /api/v1/tasks/running                  | Retrieve a list of running tasks |
| GET    | /api/v1/schedules                      | Retrieve a list of schedules |
| POST   | /api/v1/schedules                      | Create a new schedule |
| DELETE | /api/v1/schedules/delete/{schedule_id} | Delete a given schedule |

### CLI

Not written yet.  This will basically interact with the REST API though so getting it written shouldn't be difficult.

## Contributing

* Fork the repo
* Create a new branch
* Commit changes
* Write tests if applicable
* Create Pull Request

## TODO

- [ ] CLI client
- [ ] Endpoint authentication
- [ ] Proper web UI: Polymer, Bootstrap, Angular?
