/*
package clawbot provides a small Go client for the ilink Weixin QR-login flow,
message listening, replies, typing state, and media transfer used by the
OpenClaw Weixin plugin.

Minimal usage:

	client := clawbot.NewClient(clawbot.Options{})
	session, err := client.StartLogin(ctx, "")
	if err != nil {
	    log.Fatal(err)
	}

	// Show session.QRContent in your page, then wait for the user to confirm.
	account, err := client.WaitLogin(ctx, session, clawbot.WaitOptions{SaveDir: ".weixin-accounts"})
	if err != nil {
	    log.Fatal(err)
	}

The returned account contains the ilink bot token and account identifiers that
can be attached back to the same client with UseAccount for all bot operations.
*/
package clawbot
