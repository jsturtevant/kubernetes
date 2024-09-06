See [e2e-node-tests](https://git.k8s.io/community/contributors/devel/sig-node/e2e-node-tests.md)

# windows instructions
```powershell
go build -o /_output/kubelet.exe .\cmd\kubelet\ 
## needs some work to clean up the child processes, in mean time
get-process *test_e2e* | kill
get-process *kubelet* | kill
go test --container-runtime-endpoint "npipe://./pipe/containerd-containerd" --prepull-images=false --ginkgo.focus "when creating a static pod" --k8s-bin-dir "./_output/"
```

some debugging:
```powershell
kubectl --kubeconfig=.\kubeconfig get pods  -A
Get-Content .\kubelet.log -Tail 10 -Wait
Get-Content .\services.log --tail 50
```

test run:
```
W0906 15:27:49.575766   11152 test_context.go:538] Unable to find in-cluster config, using default host : https://127.0.0.1:6443
  I0906 15:27:49.575766   11152 test_context.go:553] Tolerating taints "node-role.kubernetes.io/control-plane" when considering if nodes are ready
  I0906 15:27:49.575766 11152 test_context.go:561] The --provider flag is not set. Continuing as if --provider=skeleton had been used.
  I0906 15:27:49.575766   11152 feature_gate.go:387] feature gates: {map[]}
  I0906 15:27:49.576310   11152 feature_gate.go:387] feature gates: {map[]}
  I0906 15:27:49.576310   11152 e2e_node_suite_test.go:157] failed to set rlimit on max file handles: SetRLimit unsupported in this platform
Running Suite: E2eNode Suite - C:\Users\jstur\projects\kubernetes
=================================================================
Random Seed: 1725661666 - will randomize all specs

Will run 1 of 593 specs
SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS;5;14mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS4mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS[0mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS[38;5;14mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS5;14mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSmSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS0mSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS38;5;14mSSSSSSSSSSSSSSSSSSSSSSS+SSSSSSSSSSSSSSSSSSSSSSS

Ran 1 of 593 Specs in 47.576 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 592 Skipped
```