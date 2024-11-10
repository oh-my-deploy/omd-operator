curl \
    -k \
    -X POST "https://127.0.0.1:9443/mutate-omd-com-v1alpha1-program" \
    -H "Content-Type: application/json" \
    --data @"admission_validate_request.json"