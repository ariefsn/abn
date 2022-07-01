# ABN

ABN Lookup for Go. For complete documentation please visit [ABN Lookup](https://abr.business.gov.au/) or [ABN Lookup Json](https://abn.business.gov.au/json/).

## Support

1. ABN Search
2. ACN Search
3. Name Search
4. ABN Validation (ref: [ABN Format](https://abr.business.gov.au/Help/AbnFormat))

## How to

1. Install

    ```go
    go get -u github.com/ariefsn/abn
    ```

2. Import

    ```go
    import (
      "github.com/ariefsn/abn"
    )
    ```

3. Use it

    ```go
    abn := abn.NewAbn("YOUR_GUID")

    // ABN Search
    res, code, err := abn.AbnSearch("ABN_CODE")

    fmt.Println("===== ABN =====")
    fmt.Println("[err]", err)
    fmt.Println("[code]", code)
    fmt.Println("[res]", res)

    // ACN Search
    res, code, err = abn.AcnSearch("ACN_CODE")

    fmt.Println("===== ACN =====")
    fmt.Println("[err]", err)
    fmt.Println("[code]", code)
    fmt.Println("[res]", res)

    // Name Search
    resNames, code, err := abn.NameSearch("SOME_NAME", 10)

    fmt.Println("===== Name Search =====")
    fmt.Println("[err]", err)
    fmt.Println("[code]", code)
    fmt.Println("[res]", resNames)

    // ABN Validation
    err = abn.AbnValidation("ABN_CODE")
    fmt.Println("===== ABN Validation =====")
    fmt.Println("[err]", err)

    ```
