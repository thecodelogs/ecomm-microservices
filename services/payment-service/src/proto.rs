/// Re-exports the pre-generated tonic/prost code from `gen/rust/payment/`.
///
/// The generated files live outside this crate, so we include them via a
/// path attribute rather than a build.rs / OUT_DIR approach.
pub mod payment {
    // prost message types (ProcessPaymentRequest, ProcessPaymentResponse)
    // and the tonic client/server modules are all generated in payment.rs.
    include!(concat!(
        env!("CARGO_MANIFEST_DIR"),
        "/../../gen/rust/payment/payment.rs"
    ));
}
