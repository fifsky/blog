/*
package clawbot provides a small Go client for the ilink Weixin QR-login flow
used by the OpenClaw Weixin plugin.

Minimal usage:

	client := weixin.NewClient(weixin.Options{})
	session, err := client.StartLogin(ctx, "")
	if err != nil {
	    log.Fatal(err)
	}

	// Show session.QRContent in your page, then wait for the user to confirm.
	account, err := client.WaitLogin(ctx, session, weixin.WaitOptions{SaveDir: ".weixin-accounts"})
	if err != nil {
	    log.Fatal(err)
	}

The returned account contains the ilink bot token and account identifiers that
can be reused by your own application.
*/
package clawbot
