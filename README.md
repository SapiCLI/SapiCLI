# SapiCLI - Linux API Client

### One Command Installation
```shell
curl https://raw.githubusercontent.com/SapiCLI/SapiCLI/main/install.sh | bash -
```

### Usage
![Alt text](images/panelapigui.png?raw=true "Title")
![Alt text](images/movie.gif?raw=true "Title")

The script has two uses. If you want to send a request attack within a certain time, you must use the "-timer" parameter. Required explanation is available below,

List attacks:
```shell
scli show
```
List old attacks:
```shell
scli show log
```

Normal usage format:
```shell
scli <URL OR IP ADDRESS> <PORT> <TIME> <METHOD> <CONCURRENTS>
```
Cronjob/timer mode usage format:
```shell
screen scli <URL OR IP ADDRESS> <PORT> <TIME> <METHOD> <CONCURRENTS> -timer <SECONDS>
```
and quit the screen by pressing ctrl + a + d. (In this way, it will send requests in the background according to the timer you specify.)
If you want to stop the cronjob/timer, just run this command:
```shell
pkill screen
```
screen usage: https://linuxize.com/post/how-to-use-linux-screen/




stresserus/StressersIO API Client by [forky](https://t.me/yfork) & [t13r](https://github.com/ertugrulturan/)
