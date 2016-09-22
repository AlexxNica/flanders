
#!/bin/bash

for Dir in $(go list ./... | grep -v '/vendor/'); 
do
    returnval=`go test $Dir`
    echo ${returnval}
    if [[ ${returnval} = *FAIL* ]]
    then
        exit 1
    fi  

done