# Installing cosr-front on your local machine

This guide will take you through everything you need to do to have a local instance of `cosr-front` up and running.

Please note that we only support Linux and Mac OS X at the moment.



## 1. Install dependencies on your local machine

There are only 2 dependencies you need to install to get started:

- [Docker](http://docker.com) to run containers ([Installation instructions](https://docs.docker.com/engine/installation/))
- [git](http://git-scm.com) to interact with GitHub ([Installation instructions](https://help.github.com/articles/set-up-git/))

You should make sure your Docker daemon is started. Check the [Docker documentation](https://docs.docker.com/engine/installation/) if you are unsure how to start it. For instance on Ubuntu you may need to run `sudo service docker start`.



## 2. Download the code & the Docker images

To clone this repository from GitHub, go to your local workspace directory and run:

```
git clone https://github.com/commonsearch/cosr-front.git
cd cosr-front
```

Next, there are a few Docker images available that contain all the dependencies you will need to run `cosr-front`. To download them from the Docker Hub, just run:

```
make docker_pull
```

Now you have a local copy of the code and of the Docker images on your machine!



## 3. Start the Frontend server on your local machine

The Frontend searches for documents in a local Elasticsearch instance. With `make docker_pull`, you have downloaded an image called `local-elasticsearch-devindex` that contains a few hundred pages indexed so that you can see some results right away. Let's start it:

```
make start_services_devindex
```

(To stop it later, run `make stop_services`)

The next step is to start the Frontend server itself. Just run:

```
make docker_devserver
```

Congratulations! Now you can go to the server IP shown in the logs (most likely [http://192.168.99.100:9700](http://192.168.99.100:9700) on Mac or [http://127.0.0.1:9700](http://127.0.0.1:9700) on Linux) and try some queries.
