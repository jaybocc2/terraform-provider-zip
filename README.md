# using terraform-provider-zip

```
make install
echo -e "providers {\n  terraform-provider-zip\n}\n" > ~/.terraformrc
cd ${terraform_workspace} && terraform plan
```

```
data "zip_file" "myzip" {
  files {
    lambda.py = "print 'Hello World.'"
    alt-lambda.py = "print 'GoodBye World.'"
  }
}

data "zip_file" "anotherzip" {
  files {
    lambda.py = "print 'Derp'"
  }
}


output "myzip" {
  value = "${data.zip_file.myzip.generated}"
}

output "anotherzip" {
  value = "${data.zip_file.anotherzip.generated}"
}
```

# install provider to path where terraform resides

```
make install
```

# test provider

```
make test
```

# dev build dev provider and put in gobin

```
make dev
```

# setup / test / build

```
make
```

```
make all
```

# install deps

```
make deps
```
