- 安装Jenkins
```
docker run -d --name jenkins -p 8080:8080 -p 50000:50000 -v /var/jenkins_home:/var/jenkins_home --restart always --privileged jenkins/jenkins:lts
docker run -d 
--name jenkins 
-p 8080:8080 
-p 50000:50000 
-v /var/jenkins_home:/var/jenkins_home 
--restart always 
--privileged jenkins/jenkins:lts
```