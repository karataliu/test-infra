# Acs-engine based kubernetes test design

## Pipeline Summary

1. Existing pipline

    ![existing](artifacts/1.svg)

    Currently running E2E tests on
    GCE(based on [scripts](https://github.com/kubernetes/kubernetes/tree/master/cluster/gce)) and AWS (based on [kops](https://github.com/kubernetes/kops)) all support deploying a new cluster based on built binaries. Thus they call all follow the given flow here.



1. Acs-engine pipeline

    ![acs-engine](artifacts/2.svg)

    Acs-engine supports custom kubernetes packages via [customHyperkubeImage](https://github.com/Azure/acs-engine/blob/v0.16.0/docs/kubernetes/k8s-developers.md) in api model.

    In this case, the previous flow cannot statisfy here.


## Goal
The final goal is:

1. Allows deploying cluster from source code (pull request/branch) to Azure, and running E2E test against the cluster

    The cluster would be deployed to Azure, so that a valid Azure credential is needed. The credential would contain things like 'tenant id', 'service principal id', 'service principal secert', etc.

    After that, any developers could easily set up a testing cluster as long as he provides required Azure credential and required configs.

1. Integrate with [prow](https://github.com/kubernetes/test-infra/tree/master/prow) (kubernetes test infra) to trigger tests

    A vendor organization that sponsors running test on Azure could provide their credential via [CNCF](https://github.com/kubernetes/test-infra/issues/7475#issuecomment-380281811), and later prow could pick up the credential and use it for running test.

    Examples of using those credentials could be found [here](https://github.com/kubernetes/test-infra/blob/master/prow/config.yaml#L276)


## Working steps

1. Support deploying cluster with acs-engine

    Regarding kubetest, it is to support '--up', '--down' option. And this is what [PR76250](https://github.com/kubernetes/test-infra/pull/7625) targets

1. Support uploading build artifacts to azure storage and azure container registry

    Currently kubetest only supports uploading to GCS and GCR. We need to get it support as/acr.

    The artifacts in GCS is indeed public accessable, but when testing against 100+ nodes, it will cause all the nodes to download the same package from external network. Thus it's better to make the package accessable in a storage account in same Azure region. Same for the container image.

1. Support test run via prow

    Requires communication with CNCF/prow owners for how to pass crentials.

