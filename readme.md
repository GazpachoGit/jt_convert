# PMI manager
## Description
An extension for ActiveWorkspace which allow to add and edit PMIs in the AWC Rich Text Editor

## Prerequisites
A machine with installed Teamcenter Active Workspace and TC Visualization.

## Installation
1. Copy the content of the swf_client folder to the .../awc2/stage folder. Build AWC client(awbuild)
2. The service doesn't create new folders. Prepare 3 new folders: for JT, XML, DB storage. You are going to refer them in the config.
2. For windows take the .exe file from the /bin folder. run the .exe
```
pmi_service.exe -config=<path to the config file>
```
you can find example of the config file in the /config folder