# Customer Identity check

A simple Golang program that retrieve -/+ 30k rows on mysql db and compare to Tink API Data to ensure Identity of our customers before any KYC.

Dumb and slow af, no concurrency. One of the first Go program that I wrote. 

This is a "one time" application that wont be deployed like other microservices.

Built for [Moneybounce](http://moneybounce.fr).

