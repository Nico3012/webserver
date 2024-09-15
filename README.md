# install & maintenance
this can be hosted in a repository and installed via <code>go get &lt;example.com/username/repository&gt;</code><br>
when you want to change the url to this repository, use <code>go mod edit -module &lt;example.com/username/repository&gt;</code><br>
Note, this does not change the module import path of subdirectories. You have to manually change this later.<br>

# further development
when you want to change the minimum go version, use <code>go mod edit -go &lt;version&gt;</code><br>
