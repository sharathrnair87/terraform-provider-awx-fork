file=${1}
description=$(grep "Use " ${file})

line_to_insert=$(grep -n 'return &schema.Resource{' ${file} | awk -F':' '{print $1}')

if [[ "z${line_to_insert}" != "z" ]]; then
    let "line_to_insert=line_to_insert+1"
    sed -i "${line_to_insert}i Description: \"${description}\"," ${file}
else
    echo "skipping $file"
fi

