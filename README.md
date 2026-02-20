# triple-s

Simple Storage Service implemented in Go.

A simplified S3-like object storage server that provides
REST API for bucket and object management.

## Features

- Create, list and delete buckets
- Upload, retrieve and delete objects
- Metadata stored in CSV files
- XML responses compliant with S3 specification
- Graceful error handling

## Requirements

- Go 1.22.6s
- Only standard library is used

## Build

```sh
go build -o triple-s .
```

## Run

```sh
./triple-s --port 8080 --dir ./data
```
Port and directory can be assigned.

## Usage

#### 1. Create a Bucket:
- **HTTP Method:** `PUT`
- **Endpoint:** `/{BucketName}`
- **Request Body:** Empty. Additional parameters can be passed in the request headers.
- **Behavior:**
  - Validate the bucket name to ensure it meets Amazon S3 naming requirements (3-63 characters, only lowercase letters, numbers, hyphens, and periods).
  - Ensure the bucket name is unique across the entire storage system.
  - If the bucket name is valid and unique, create a new entry in the bucket metadata storage.
  - Return a `200 OK` status code and details of the created bucket, or an appropriate error message if the creation fails (e.g., `400 Bad Request` for invalid names, `409 Conflict` for duplicate names).

Rely on the [documentation](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CreateBucket.html#API_CreateBucket_Examples)

#### 2. List All Buckets:
- **HTTP Method:** `GET`
- **Endpoint:** `/`
- **Behavior:**
  - Read the bucket metadata from the storage (e.g., a CSV file).
  - Return an XML response containing a list of all matching buckets, including metadata like creation time, last modified time, etc.
  - Respond with a `200 OK` status code and the [XML list of buckets](https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html#API_ListBuckets_Examples).

#### 3. Delete a Bucket:
- **HTTP Method:** `DELETE`
- **Endpoint:** `/{BucketName}`
- **Behavior:**
  - Check if the specified bucket exists by looking it up in the bucket metadata storage.
  - Ensure the bucket is empty (no objects are stored in it) before deletion.
  - If the bucket exists and is empty, remove it from the metadata storage.
  - Return a `204 No Content` status code if the deletion is successful, or an error message in XML format if the bucket does not exist or is not empty (e.g., `404 Not Found` for a non-existent bucket, 409 Conflict for a non-empty bucket).

Don't forget to process the data and save the corresponding metadata in your CSV file.

### Ensuring Unique and Valid Bucket Names:

#### Bucket Naming Conventions:

- Bucket names must be unique across the system.
- Names should be between 3 and 63 characters long.
- Only lowercase letters, numbers, hyphens (`-`), and dots (`.`) are allowed.
- Must not be formatted as an IP address (e.g., 192.168.0.1).
- Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.

#### Validation Implementation

- Use regular expressions to enforce naming rules.
- Check the uniqueness of a bucket name by reading the existing entries from the CSV metadata file.
- If the bucket name does not meet the rules, return a `400 Bad Request` response with a relevant error message.



### API Endpoints for Object Operations

#### 1. Upload a New Object:
- **HTTP Method:** `PUT`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Request Body:** Binary data of the object (file content).
- **Headers:**
  - `Content-Type`: The object's data type.
  - `Content-Length`: The length of the content in bytes.
- **Behavior:**
  - Verify if the specified bucket `{BucketName}` exists by reading from the bucket metadata.
  - Validate the object key `{ObjectKey}`.
  - Save the object content to a file in a directory named after the bucket (`data/{BucketName}/`).
  - Store object metadata in a CSV file (`data/{BucketName}/objects.csv`).
  - Respond with a 200 status code or an appropriate error message if the upload fails.
  - If an object with the same name already exists, it must be overwritten.


#### 2. Retrieve an Object:
- **HTTP Method:** `GET`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Behavior:**
  - Verify if the bucket `{BucketName}` exists.
  - Check if the object `{ObjectKey}` exists.
  - Return the object data or an error.


#### 3. Delete an Object:
- **HTTP Method:** `DELETE`
- **Endpoint:** `/{BucketName}/{ObjectKey}`
- **Behavior:**
  - Verify if the bucket and object exist.
  - Delete the object and update metadata.
  - Respond with a `204 No Content` status code or an appropriate error message.

## Examples:

Firstly, run the server with default port and directory.

### Bucket creation.
```sh
curl -i -X PUT http://localhost:8080/my-bucket
```

#### Output
```sh
HTTP/1.1 409 Conflict
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 20 Feb 2026 06:18:56 GMT
Content-Length: 141

The requested bucket name is not available.The bucket namespace is shared by all users of the system. Select a different name and try again.
```


### Bucket deletion 
```sh
curl -i -X DELETE http://localhost:8080/my-bucket
```

#### Note
"Non-empty Bucket" is being output when:
- the directory does not exist but it may have a record in bucket.csv
- the directory is non-empty

#### Output
```sh
HTTP/1.1 204 No Content
Date: 2026-02-20T07:02:22Z
```

### List of all buckets

Beforehand, add put buckets.

```sh
curl -X GET http://localhost:8080/
```

#### Output
```sh
<ListAllMyBucketsResult>
  <Buckets>
    <Bucket>
      <Name>my-bucket</Name>
      <CreationDate>2026-02-20T11:44:36+05:00</CreationDate>
      <LastModifiedTime>2026-02-20T12:02:22+05:00</LastModifiedTime>
      <Status>deleted</Status>
    </Bucket>
    <Bucket>
      <Name>my-cucket</Name>
      <CreationDate>2026-02-20T12:06:08+05:00</CreationDate>
      <LastModifiedTime>2026-02-20T12:06:08+05:00</LastModifiedTime>
      <Status>active</Status>
    </Bucket>
    <Bucket>
      <Name>my-cucket2</Name>
      <CreationDate>2026-02-20T12:06:10+05:00</CreationDate>
      <LastModifiedTime>2026-02-20T12:06:10+05:00</LastModifiedTime>
      <Status>active</Status>
    </Bucket>
    <Bucket>
      <Name>my-cucket3</Name>
      <CreationDate>2026-02-20T12:06:12+05:00</CreationDate>
      <LastModifiedTime>2026-02-20T12:06:12+05:00</LastModifiedTime>
      <Status>active</Status>
    </Bucket>
    <Bucket>
      <Name>my-cucket4</Name>
      <CreationDate>2026-02-20T12:06:14+05:00</CreationDate>
      <LastModifiedTime>2026-02-20T12:06:14+05:00</LastModifiedTime>
      <Status>active</Status>
    </Bucket>
  </Buckets>
```

### Put object

```sh
curl -i -X PUT -F "file=@test.jpg"   http://localhost:8080/my-cucket/test.jpg -v
```

#### Output
```sh
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> PUT /my-cucket/test.jpg HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.5.0
> Accept: */*
> Content-Length: 4799857
> Content-Type: multipart/form-data; boundary=------------------------W8vKdmvYLucZQaNYHR0c2A
> Expect: 100-continue
> 
* Done waiting for 100-continue
< HTTP/1.1 100 Continue
HTTP/1.1 100 Continue

* We are completely uploaded and fine
< HTTP/1.1 200 OK
HTTP/1.1 200 OK
< Content-Length: 4799857
Content-Length: 4799857
< Content-Type: multipart/form-data; boundary=------------------------W8vKdmvYLucZQaNYHR0c2A
Content-Type: multipart/form-data; boundary=------------------------W8vKdmvYLucZQaNYHR0c2A
< Date: Fri, 20 Feb 2026 08:30:55 GMT
Date: Fri, 20 Feb 2026 08:30:55 GMT

< 
* transfer closed with 4799857 bytes remaining to read
* Closing connection
curl: (18) transfer closed with 4799857 bytes remaining to read
```

#### Note
Also works with
```sh
curl -X PUT -H "Content-Type: image/jpeg" --data-binary @test.jpg \
  http://localhost:8080/my-cucket/test.jpg -v
```
The last message "curl: (18) transfer closed with 4799857 bytes remaining to read" will be gone if Content-Length will be assigned.

### Download file 

```sh
curl http://localhost:8080/my-cucket/test.jpg -o downloaded.jpg
ls -lh downloaded.jpg
```

#### Output

```sh
-rwxrwxrwx 1 user user 4.6M Feb 20 14:11 downloaded.jpg
``` 

___

```sh
curl -i http://localhost:8080/my-cucket/test.jpg
```

#### Output
```sh
HTTP/1.1 200 OK
Content-Length: 4799659
Content-Type: image/jpeg
Lastmodified: 2026-02-20T13:41:12+05:00
Date: Fri, 20 Feb 2026 09:15:48 GMT

```

### Delete object

```sh
curl -i -X DELETE http://localhost:8080/my-cucket/test.jpg
```

#### Output

```sh
HTTP/1.1 204 No Content
Content-Type: None
Lastmodified: 2026-02-20T14:48:52+05:00
Date: Fri, 20 Feb 2026 09:48:52 GMT

```

As you check the directory you will see that previously uploaded file is absent now as well as the record of it in object.csv

## Interesting cases

Rerun the server and clean data directories. Before the next tests create "my-bucket".

### Put Object with empty object-key

```sh
curl -i -X PUT \
  --data-binary @test.jpg \
  http://localhost:8080/my-bucket/
```
It just tries to create bucket instead of upload the given binary.

### Invalid object-key 

```sh
curl -i -X PUT \
  --data-binary @test.jpg \
  http://localhost:8080/my-bucket/test?.jpg
```

It normalizes URL and remove everything that follows '?'

```sh
curl -i -X PUT \
  --data-binary @test.jpg \
  http://localhost:8080/my-bucket/../hack.jpg
```

It implements traversing in file system so API URL become 
```sh
/hack.jpg
```
which creates buckey 'hack.jpg'