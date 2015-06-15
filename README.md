#s3cp
An scp-style command-line tool for uploading and downloading files to and from an AWS S3 bucket.

### License
Released under the Apache 2.0 license

### Current Status
s3cp is currently in beta verion 0.1
* Uploading to an S3 bucket works
* Errors are trapped and reported properly
* Exit codes work as expected (0 for success, non-0 for failure)
* *Downloads are not working yet*

# How to use s3cp

### Step 1 - Install s3cp
You can download the executable from the homepage, or you can get the sources and compile. The
executable is statically linked (as all Go programs are), so you can download the executable
and run it without installing Go.

If you intend to compile, you'll need to install the Go language
(for more information about installing Go on your platform, visit the [Go homepage](http://golang.org))

If you have Go installed, you can get s3cp the "Go way" like this:
```
go install http://github.com/tdunnington/s3cp
```

### Step 2 - Setup your AWS S3 Bucket
Setting up an S3 bucket to work with s3cp is really about security. The policy editors have become much
easier to use over the last year.

Please note that these are my recommended settings. You are ultimately responsible for the security of your bucket!

_Setting up with a user account_
This is the only method you can use if your server is outside of AWS. If your server is inside AWS (hosted via
EC2), then the "Setting up with EC2 IAM roles" is a better choice for you, but if you did not create an IAM
role when you created the server, you will either have to use the user account method, or re-launch your server
with an IAM role.

In this method, you will create an IAM user and an IAM group in AWS. The group will contain the policy permissions
that will allow you to upload and download to the target bucket. The advantage of this method is that you can easily
swap out credentials, but the downside is that the credentials have to be stored on the system, where they are
vulnerable.

Create a new group in IAM (let's call it "s3cp_test"). Do not create a policy for the group when you create it; 
we'll add the policy in the next step:

...TODO add images...

Now add an inline policy to the group:

...TODO add images...

Use the "Custom Policy" option and paste this sample policy into the editor. Note that you have to change the 
bucket name:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "removed for security",
            "Effect": "Allow",
            "Action": [
                "s3:DeleteObject",
                "s3:DeleteObjectVersion",
                "s3:GetBucketPolicy",
                "s3:GetObject",
                "s3:GetObjectAcl",
                "s3:GetObjectVersion",
                "s3:GetObjectVersionAcl",
                "s3:ListBucket",
                "s3:ListBucketMultipartUploads",
                "s3:ListBucketVersions",
                "s3:ListMultipartUploadParts",
                "s3:PutObject",
                "s3:PutObjectAcl",
                "s3:PutObjectVersionAcl"
            ],
            "Resource": [
                "arn:aws:s3:::<PUT BUCKET NAME HERE>/*"
            ]
        }
    ]
}
```

_Note_: I am not certain that these are the minimal permissions necessary. I have not found documentation that 
makes it clear what permissions are required to add an item to a bucket using the [s3manager](http://github.com/aws/aws-sdk-go/service/s3/s3manager) that is used by s3cp.

Now, create a user for s3cp, and save the AWS security keys. These screenshots show an actual AWS key pair;
rest assured that key pair has been deleted and is not valid.

...TODO add screenshots...

Hold on to the key pairs; you will need them for the next step. Go on to Step 3 now.

_Setting up with EC2 IAM roles_
This method allows you to use AWS resources from an EC2 instance without having to store the key pair on the
server anywhere; AWS will know the server has permissions based on the server role.

In order for this method to work, the IAM role must be applied to the server when it is created.
TODO

### Step 3 - Create your credentials file/environment variables
TODO

### Step 4 - Use s3cp
TODO
