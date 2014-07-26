# Mailrouter

A dynamically configurable mail router with a web interface, written in Go. It requires Go 1.2 or newer.

It is intended to manage and redirect mail sent by deployed applications without having to alter the application or the production mail server configuration.

## Features

* Define filters (routing rules) on From address, To address, Subject header and originating IP.
* Ordering of filters.
* The ability to readdress mail matching a filter.
* A web interface for configuring SMTP routes and routing rules (called filters).
* A customisable listening address and port for both HTTP and SMTP interfaces.
* Logging of delivered and dropped mail messages.
* The ability to set a default route for mail.
* A human-readable configuration file in JSON format.
* IPV6 support, including theoretical use as a gateway for IPV6-only servers to route mail to an IPV4-only mail server.
* Routing loops are allowed.

## Intended Use

Mailrouter is intended to be placed between mail-sending applications and proper SMTP servers. It should then be used to route, forward or drop mail messages as defined by filters on a dynamic basis.

For example: An organisation has a production mail server mail.example.com, a deployment of the [Mailcatcher](http://mailcatcher.me/) software at mc.example.com, and two versions of a deployed web application that can send user-defined emails. Two development versions of the application run on 10.0.0.1 and 10.0.0.2. A production version runs on 10.0.1.1.

The organisation could set up:

* A Route named Outbound to mail.example.com. All mail sent to this Route will be delivered to mail.example.com.
* A Route named QA to mail.example.com with the To field set to qa@example.com. All mail sent to this Route will be delivered to qa@example.com on the mail server, acting as an automatic forward. This will allow the QA team to examine mail messages as they would be received by regular users.
* A route named Mailcatcher to mc.example.com.

It could then set up the following:

* A Filter named Development with the Originating IP of 10.0.0.0/24 to route mail to the Mailcatcher Route. This would ensure no mail sent from the 10.0.0.1 server, or any other server in the 10.0.0.x IP range, would escape to the public internet.
* A Filter named Test Mail with the Originating IP of 10.0.1.1 and the To address of "test" to route mail to the QA Route. This would ensure any accounts on the production server with a To address containing the word "test", in either the user or domain section, would be redirected to the QA team.
* The Outbound route as the default route.

Soon, the application gains new users. However, some of these users are unpleasant, and the organisation adds a Filter with a Subject containing the word "badsite.com" to route mail to the Drop Route. This would drop all mail with links to badsite.com because email from those accounts is to be ignored.

Later, when the application is no longer under active development, the Route for the Test Mail filter can be set to Drop, as the QA team has been reassigned to the next big thing and no longer needs testing messages from the application.

## Compiling

First get the dependencies:

	go get github.com/mhale/smtpd
	go get github.com/streadway/simpleuuid
	go get github.com/jteeuwen/go-bindata

Then install each of them with:

	go install

## Command Line Usage

There are four command line options.

* -h prints the help message
* -conf specifies the path to store the configuration file. The default is /etc/mailrouter.conf.
* -http specifies an address & port for the HTTP server to listen on. The default is all addresses and port 8080.
* -smtp specifies an address & port for the SMTP server to listen on. The default is all addresses and port 2525.

The default values are equivalent to:

	mailrouter -conf=/etc/mailrouter.conf -http=:8080 -smtp=:2525

To bind to a specified IPv4 address:

	mailrouter -http=127.0.0.1:80 -smtp=127.0.0.1:25

Any IP:port format accepted by Go will work, however IPv6 addresses have not been tested yet.

## Tips

* Create Routes first, so the drop-down Route selector is populated when Filters are created.
* Define Filters in order beginning at 100, numbering the second Filter as 200, the third as 300, and so on. This provides flexibility later when inserting new Filters between existing Filters.
* Filter fields are logical AND operations i.e. they must all match for the Filter to match. Place more specific Filters before general Filters.
* Filters will be checked in the order displayed on the Filters page.
* If no routes are configured, all mail will be dropped. This can be useful when your application requires a mail gateway but you don't care about the mail.

## To Do

* SMTP authentication support.
* SSL/TLS support.
* Verify that use as an IPv4 to IPv6 bridge works.
* Hostname support in the Originating IP field.
* Mail header matching and overriding.
* Filtering by body text.
* Filtering by attachments: file count, file size, MIME type, etc.
* Full end-to-end testing. Currently only basic testing of the filtering functionality is implemented.

## Development

Pull requests are welcome. To edit the assets or views, you will need to get the [go-bindata](https://github.com/jteeuwen/go-bindata/) command line tool (which comes with the package) to generate the bindata.go file:

	cd $GOPATH/src/github.com/jteeuwen/go-bindata/go-bindata
	go install

Then generate a shim with:

	go-bindata -debug assets views

This will allow you to update the views and hit Refresh in your browser to see your changes without recompiling or restarting the program.

## Bugs / Issues

Please report any bugs you find as Github Issues.

## Licensing

Copyright © 2014 Mark Hale. This program is released under the MIT License. See the LICENSE file for details.
