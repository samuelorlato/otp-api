# OTP API
A Go API that sends verification code (OTP)

## Routes

<details>
 <summary><code>POST</code> <code><b>/email/otp</b></code> <code>(sends OTP to email)</code></summary>

##### JSON Body Params

> | name   | type     | data type |
> | ------ | -------- | --------- |
> | email  | required | string    |

![/vote](./assets/send.png)

</details>

<details>
 <summary><code>POST</code> <code><b>/verify/otp</b></code> <code>(verifies email)</code></summary>

##### JSON Body Params

> | name                   | type     | data type |
> | ------                 | -------- | --------- |
> | OTP                    | required | string    |
> | encryptedVerification  | required | string    |
> | email                  | required | string    |

![/vote](./assets/verify.png)

</details>