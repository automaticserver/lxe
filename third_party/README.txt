These packages are copied from <https://github.com/kubernetes/kubernetes>.

The kubernetes main repo does not correctly set the go modules, see bug
<https://github.com/kubernetes/kubernetes/issues/79384. The workaround we used
before was setting the there "suggested" replace directives. Since Go 1.18 there
is a convenient `install` command, but it does not allow installation from
repos with replace directives. So we had to copy the used modules and placed
here. CRI-O and probably the other CRI shims copied/used to copy those modules
as well. E.g.
  https://github.com/cri-o/cri-o/tree/8d4d158935929800c4300b82eb4b5a83ded400f4/server/streaming
