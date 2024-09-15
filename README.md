# install & maintenance
this can be hosted in a repository and installed via <code>go get &lt;example.com/username/repository&gt;</code><br>
when you want to change the url to this repository, use <code>go mod edit -module &lt;example.com/username/repository&gt;</code><br>
Note, this does not change the module import path of subdirectories. You have to manually change this later.<br>

# disable go proxy
To list the current value of all environment variables in go use: <code>go env</code>
To modify a environment variable use <code>go env -w &lt;VARIABLE&gt;=&lt;value&gt;</code>
To fetch packages directly from your server, you can set the GONOPROXY to your dependency domains: <code>go env -w GONOPROXY=&lt;example.com&gt;</code><br>
To fetch packages directly from your server without a sum database, you can set the GONOSUMDB to your dependency domains: <code>go env -w GONOSUMDB=&lt;example.com&gt;</code><br>

# todo
read more about go get from proxy server because go get command requests the packages from a internal proxy server, that requires a internet connection.

# further development
when you want to change the minimum go version, use <code>go mod edit -go &lt;version&gt;</code><br>

# naming convention
this module should be named "github.com/Nico3012/webserver" but when using go workspaces, the go mod tidy command tries to download the package because it not recognizes the go.work file. by simply using "webserver" as name, the go mod tidy command seas, this is not a url and does not fetch the package.
