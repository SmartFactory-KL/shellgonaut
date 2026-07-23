# General ideas for this

This package should be a Shell Client that could be used to interface with
 - Discovery - probably on its own, as it differs from the rest quite a bit
 - Registry - maybe using RepoClients in a cache, unsure wether or not to have own RegistryClient
 - Repository - already there, although not completed yet
 - Services - this is up for discussion, as it would then need to be added to the already existing repo clients

In addition to the prior things it should also offer quality of life improvements for
 - traversing shells and submodels, especially reducing "all the nil checks" and stuff
 - reading and writing of simpletons, like Properties
 - copy, create
 - reading with the idShortPath

But that would be separate from the pure clients. The client should always return a "naked" types.IClass.

One caveat might be the golang aas core works, as that is only a release candidate as of now.

For Registry Use there might be the need for a new type of client, because 
a simple RegistryClient would not be enough (it would only cover the registry api)
A name will be required for that new layer.
Like: RegistryBasedRequester or something, showing that this is using other clients in combination to get
to the shell. In the end it would probably have methods like "GetShell(shellID)" and it would do all the registry stuff in the background for you
