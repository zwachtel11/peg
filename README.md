# peg, the minimalist kubernetes package manager
Kubernetes is hard ... helm makes it harder. Try peg for to "peg" up a seemless Kubernetes experience.

peg is a oras based storage for kubernetes manifests. a minimalist package manager geared at helping the deployment story for private container registries.

With peg a user can seemlessly push their kubernetes application manifest into a OCI compliant container registry. Then when they need to deploy their appliacations, they can pull down their manifests or seemlessly deploy their applications to the target kubernetes cluster.

## install cli

* Install from the latest [release artifacts](https://github.com/zwachtel11/peg/releases):

  * Linux

    ```sh
    curl -LO https://github.com/zwachtel11/peg/releases/download/v0.0.1/peg
    mv peg /usr/local/bin/
    ```

  * macOS

    ```sh
    curl -LO https://github.com/zwachtel11/peg/releases/download/v0.0.1/peg-darwin
    mv peg-darwin /usr/local/bin/peg
    ```

  * Windows

    Add `%USERPROFILE%\bin\` to your `PATH` environment variable so that `peg.exe` can be found.

    ```sh
    curl.exe -sLO  https://github.com/zwachtel11/peg/releases/download/v0.0.1/peg.exe
    copy peg.exe %USERPROFILE%\bin\
    set PATH=%USERPROFILE%\bin\;%PATH%


## usage

Lets say you have a custom redis configuration that your team commonly uses on clusters.

For instance here is `pods/config/redis-pod.yaml`

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: redis
spec:
  containers:
  - name: redis
    image: redis:5.0.4
    command:
      - redis-server
      - "/redis-master/redis.conf"
    env:
    - name: MASTER
      value: "true"
    ports:
    - containerPort: 6379
    resources:
      limits:
        cpu: "0.1"
    volumeMounts:
    - mountPath: /redis-master-data
      name: data
    - mountPath: /redis-master
      name: config
  volumes:
    - name: data
      emptyDir: {}
    - name: config
      configMap:
        name: example-redis-config
        items:
        - key: redis-config
          path: redis.conf
```

Instead of copy and pasting this config from machine to machine or cloning git repos you decide to use peg to seemlessly push and pull your kubernetes manifests.

First use peg to login into your private OCI compliant container registry

```sh
peg login pegimages.azurecr.io -u <CR_USERNAME> -p <CR_PASSWORD>
```

Once you are authenticated, you are ready to push manifests into the registry.

```sh
peg push --file pods/config/redis-pod.yaml --manifest pegimages.azurecr.io/myteam-redis:latest
```

For here on your production machines you should be able to pull your config just as you expect for a docker image.

```sh
peg pull --manifest pegimages.azurecr.io/myteam-redis:latest --outfile myteam-redis.yaml
```

Or if you already know the config is good to go, cut out the middle step and use peg to deploy your application straight to your cluster.

```sh
peg deploy --manifest pegimages.azurecr.io/myteam-redis:latest --kubeconfig=kubeconfig
```

Or if your kubeconfig is already set in `~/.kube/config`

```sh
peg deploy --manifest pegimages.azurecr.io/myteam-redis:latest
```

## Roadmap

* Docker Secrets embedded in the deployment yamls for seamless private container pull.
* "pegify" the helm chart experience with a simple push & pull flow for helm charts.
* Post an issue!
