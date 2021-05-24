package alipayment

const (
	AppId           = "2021002134605526"
	AliPayPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuddg5Uz1RKyGz3kPjgLWh8Mcp8DUlC1SN6N1mzf4ZyusGLJPu6ejoUE8A97kiPugILjkIks9gRq9Sc8yl1aKh7w1RmUyj8IyYKlr9ykbpdmXZ1KJnFG6bRbK2jUHMqf0k+xBv4xy+E0qE/UgueBLxtDtaBBynCuYI61AVsD2D1haqqIUITWKvfHjcMd5eEQ8KlsHGSfHbwOd+JsCGSs8VH0wkcChR4+56dtzhBWN0fdZiPZnUlL20yA/QsJloGmiCwCR351XsGzWn1ElA97iPTCajFvfVkLWciW0XggcG3oEPBLdVvArF1MAzCL8T7ehVP6V7Crn01tzLQndrsDo4QIDAQAB"
	PrivateKey      = "MIIEpAIBAAKCAQEAuddg5Uz1RKyGz3kPjgLWh8Mcp8DUlC1SN6N1mzf4ZyusGLJPu6ejoUE8A97kiPugILjkIks9gRq9Sc8yl1aKh7w1RmUyj8IyYKlr9ykbpdmXZ1KJnFG6bRbK2jUHMqf0k+xBv4xy+E0qE/UgueBLxtDtaBBynCuYI61AVsD2D1haqqIUITWKvfHjcMd5eEQ8KlsHGSfHbwOd+JsCGSs8VH0wkcChR4+56dtzhBWN0fdZiPZnUlL20yA/QsJloGmiCwCR351XsGzWn1ElA97iPTCajFvfVkLWciW0XggcG3oEPBLdVvArF1MAzCL8T7ehVP6V7Crn01tzLQndrsDo4QIDAQABAoIBABVpgmmJD21lL3cyLJ+hJhSoX2HYruWPvzxX4hILRdylRIPn75Xgf9LxkDel89DwpQeAJhDpeGGqsgqSC/Mx7m4rOjwtJsE3j9RhUOY+X0ghpPcrpF1VcVRmXuL7XyPwBBcTDpRCW3DOez2nbajc9f4nTmNCGeJyh2n86T/VpcQoqBLJo73ch/qngtxMdYTXmO7VVvG/KlPlf73xriEzLnyYZcNGg7LUCVqrNPEvTkezGGqfDvF4ta+e2yNIM+O9/RcV0tWkHq8ghKC7UahCldoTAKImrqJ/lHmc+1TIRhpEpthQ+VlgeqGsyO4/PEsLW4VnUQdIvFwuywoqWmpDev0CgYEA/sfzLzg8Igs8vuoP+1BPFJW3va9YDmVZk1xsG8kKvB4MvX6TFxwRAvX0ykXr0q1lZmozv3L+f73ergjW5j9nt5/D8qGV6I3z9pjs1SPMESp+yn918+LFRmUZCFJs2AkKqbcU90A+iwgewC0PI71kJLzKAqjNnRtmU6mCDu+M6VcCgYEAurr+JBZ/F9VoTS6gYJUnk2B8cROgdbRnkgC5feV9+0Dinzo5eXG3djVZKf4LX0zMzQ6SmKp8vIE32GrdL01qcZY6fSd8N91fFVh46WbC1hF6oPzz/yR9XDUBYmHvuL8+78K8I4I0l2ixgLJUJCpS/i3bmYPJqEJX8TQ38mRHhIcCgYEA3kDqoRx7SG9pFFn2ixLpiVl5qZzWQeijlWsgW7SrCqNkFcHe6l/vbxfD7NfF+kiMliS6852K6UBmiobngH97KLHUK0pODUoXsdyBbQUHNUpOxLMf4BrIYLo8dggBLwvkI1y0i2Odq4bv3FDyTgke3PVbe6ppg02tm/nYnNLXp40CgYBjsd5qChDE+s17F7/MkRXWG5eVD7OF5FG/o99lFhfTA2t7M6zn6AzZcdYvE1GjCe/2MyeVlsCOLDdStots69o+1sRXebNaaDiRvvCfPJwdiqt8NFFZEXRUvIAGtN82NfU2MTQdiTDm/aB3Y76EqIwy5Ozv42rRCMrHHugX4+5DFQKBgQCn0+Pd6PtQtrUzru3YzdIlaieIiuG4e7CBm4zLysidSxDc3TcUS9IW77ZraOZfq++UoxjaeX7ktA4eOMx2YyfvKdnLFm8j346ys0Vez3TpOr7vn+uzH3+ze/zE7pKzITAmgX7tO0qNWxnSS4SOVVK9/tjuEHupSafB8A3fIAQS0Q=="
)

const (
	AppCertPublicKey     = "appCertPublicKey_2021002134605526.crt"
	AliPayRootCert       = "alipayRootCert.crt"
	AppCertPublicKeyRSA2 = "alipayCertPublicKey_RSA2.crt"
)

const (
	NotifyUrl       = "120.24.151.34:80/pay/ali/callback"
	RefundNotifyUrl = "https://api.fmg.net.cn/"

	UnifiedApplyNotifyUrl = "https://api.fmg.net.cn/pay/"
	RefundApplyNotifyUrl  = "https://api.fmg.net.cn/pay/"
)