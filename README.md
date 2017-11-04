# SimpleAzureUploader-go
Simple file uploader for azure storage for simple cases.
Written in golang so can run without deps.


## Usage:
azure-uploader-{OS} -accountname=<> -accountkey=<> -containername=<> -filename=<>

or you can replace the flags with the relevant environment variables:
ACCOUNT_NAME, ACCOUNT_KEY, CONTAINER_NAME.

on the CLI you can also specify the name of the target blob with -targetname.

and the content type with -contenttype.


