# Testing

## Unit tests

Unit tests can be run with `make test` or `go test ./...` and don't require any environmental setup. All external interfaces are tested against mocks.

## cri-tools' critest

LXE has a `--critest` option which prepares some default images to use in the `critest` command since LXE uses LXD and that doesn't support OCI images. It also rewrites *some* image requests as the hardcoded image names in critest are just not available in LXD. If it does so it clearly writes this in the output like `CRITest: %s`. This feature is considered experimental as the images might still not meet the wished requirments.

LXE will never pass all the critest cases due to the fact that LXD does not support OCI images and LXD might not offer all the runtime features docker provides. LXD does **system containers**, while docker and similiar OCI runtimes do **application containers**.

Run lxe in critest mode:

```
$ lxe --log-level info --critest
WARNING[06-08|13:44:35.631] starting lxe...                               builddate="Wed, 08 Jun 2022 13:44:34 +0200" buildnumber=undef gitcommit=3a8f4f75b2fe8d0d7742e66175a529e2ed67fee3 gittreestate=dirty packagename=lxe version=0.0.0
INFO   [06-08|13:44:35.634] Connected to LXD                              lxdsocket=/var/lib/lxd/unix.socket
WARNING[06-08|13:44:35.634] CRITest mode enabled                         
INFO   [06-08|13:44:35.638] downloading default image                    
INFO   [06-08|13:44:38.713] default image ready: 'local:critest/default' 
INFO   [06-08|13:44:38.727] creating webserver image                     
INFO   [06-08|13:44:55.801] webserver image ready: 'local:critest/webserver' 
WARNING[06-08|13:44:55.802] CRITest ready: Use --test-images-file=/tmp/lxe-critest-images-file.yaml in your critest command 
INFO   [06-08|13:44:56.122] started lxe CRI shim                          socket=/tmp/lxe.sock
INFO   [06-08|13:44:56.123] started streaming server                      baseurl="http://localhost:44124" endpoint="localhost:44124"
```

And use the prepared test images file as provided in the output in your critest command.

For the image specific tests:

```
$ critest --runtime-endpoint unix:///tmp/lxe.sock --test-images-file=/tmp/lxe-critest-images-file.yaml --ginkgo.focus "image" --ginkgo.skip "should not have Uid|Username empty" --ginkgo.slow-spec-threshold "30s"
critest version: v1.24.1
Running Suite: CRI validation - /root
==============================================
Random Seed: 1654688963

Will run 6 of 79 specs
SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS••••S••SSSSSSSSSSSSSSS

Ran 6 of 79 Specs in 16.681 seconds
SUCCESS! -- 6 Passed | 0 Failed | 0 Pending | 73 Skipped
[...]

PASS
```

Or remove the focus and skips and just let it run over the full test suite:

```
$ critest --runtime-endpoint unix:///tmp/lxe.sock --test-images-file=/tmp/lxe-critest-images-file.yaml --ginkgo.slow-spec-threshold "30s"
[...]
```
