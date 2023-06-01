# kubeconfig-splitter
Split files when there are multiple contexts in kubeconfig

# Usage
1. First, clone this repository.
2. Then place the kubeconfig file you want to split in the repository.
3. Run the tool.
  ```shell
  go build main.go
  ./main
  ```
4. A number of split kubeconfig files are created in the /tmp directory for the number of contexts.<br>
   At that time, a file in the form of kubeconfigN is created in the temp directory of the OS.
