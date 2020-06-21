provider aws {
  region  = "us-east-1"
  version = "~> 2.51"
}

provider aws {
  alias  = "config_registry"
  region = "us-east-1"
}

terraform {
  required_version = "~> 0.12.21"
  backend "s3" {
    bucket  = "spicy-omelet-terraform"
    key     = "circleci/hi.tfstate"
    // Use Server-Side Encryption with Amazon S3-Managed Keys (SSE-S3)
    encrypt = "true"
    region  = "us-east-1"
  }
}