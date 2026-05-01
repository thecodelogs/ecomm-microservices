# Go (user service)
buf generate --template buf.gen.go.yaml --path proto/user/v1
# notification service
buf generate --template buf.gen.go.yaml --path proto/notification/v1
# product service
buf generate --template buf.gen.go.yaml --path proto/product/v1
# shipping service
buf generate --template buf.gen.go.yaml --path proto/shipping/v1
# cart service
buf generate --template buf.gen.go.yaml --path proto/cart/v1
# order service
buf generate --template buf.gen.go.yaml --path proto/order/v1

# Rust (payment service)
buf generate --template buf.gen.rust.yaml --path proto/payment/v1
