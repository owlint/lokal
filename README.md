# Lokal - Easily run your kubernetes applications in local.
Lokal is a command-line tool that simplifies the process of local development alongside a kubernetes cluster. With Lokal, you can seamlessly inject remote variables and secrets into your local applications.

## Usage

### Run a local application
```
$ lokal run
```
#### Options
   - `--namespace value`:    Namespace within cluster. Overrides config file if provided.
   - `--deployment value`:   Deployment within namespace. Overrides config file if provided.
   - `--container value`:    Container within deployment. Will use pod name if ommited.
   - `--command value`:      Command to run.
   - `--force-namespace`:    Append namespace to url environnement variables referencing local deployment. Ex: http://myapp/ will be converted to http://myapp.mynamespace/. This feature is useful when running lokal alongside Telepresence. (default: true)
   - `--config value`:       Path to the local config file. (default: "./lokal.yaml")
   - `--kube-config`: value  Path to the kube config file. (default: "~/.kube/config")
   - `--help, -h`: show help
   
#### Configuration File

You can also configure theses parameters using a YAML configuration file containing the following fields:

- `namespace` (string): The Kubernetes namespace where the deployment is located.
- `deployment` (string): The name of the deployment in which the container resides.
- `container` (string): The name of the container from which to retrieve the environment variables.
- `command` (string): The command to execute.
- `env` (list): A list of environment variables to inject. Each environment variable is defined with a `name` and `value`.

##### Exemple
```
namespace: such-namespace
deployment: such-deployment
container: such-container
command: "env"
env:
  - name: HELLO
    value: World
```