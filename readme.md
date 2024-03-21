## UntisChangeNotifier

UntisChangeNotifier is a project that uses the untis API to notify its user in case changes to the timetable and new
homeworks are detected. It uses the unofficial [UntisAPI](https://github.com/tomroth04/untisAPI) that I wrote myself.
The project stores the last state of the timetable and the homeworks in a Key-value database on disk.

The project continously checks for changes and notifies the user through a notification service.

The project supports multiple notification services
through [containrrr/shoutrrr](https://github.com/containrrr/shoutrrr).


## Usage

The project is designed to be run as a docker container.
In the future I might also provide a standalone binary. In the meantime feel free to compile it yourself.
For a guide on the configuration check out the [Wiki](https://github.com/tomroth04/UntisChangeNotifier/wiki)


### Information
The bot is configured to pull to have a look at the timetable more often in the morning when changes usually happen to
the timetable at my school. This could be changed in the future to be more dynamic.
This is to reduce the load on the untis servers.