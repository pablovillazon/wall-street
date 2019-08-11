Integration Patterns

FINAL PROJECT

1. Analysis - Patterns Proposed
2. C4 Mode


1. Analysis - Patterns Proposed
Below the problem description.
2 Wall Street
A Wall Street firm handles trading of many different commodities in national and international levels, all the transactions are managed internally providing a APIs for external consumers to fetch reports and data about their stocks, transactions and trade forecasting. They are running on an average of 2000 transactions per second on weekdays and 100 transactions per second on weekends.
Two months ago, the firm has been merged with two other firms, this means that their customer base has doubled overnight. Because of it they have had latency issues and outages, increasing the exchange rate from 200 concurrent requests per second to 1000 concurrent requests in their main servers. The increased load is causing delays in the internal transactions and overall dissatisfaction in all traders.
What they can do to improve/scale up their service quality without affecting transactions and external customers?

Current Scenario
The initial design has and intermediate service working as the API endpoint serving the customer requests. Fig 1.
![Alt text](http://jkin.be/content/WallStreet-Initial.png)
<details> 
<summary></summary>
custom_mark10
  digraph G {
    size ="4,4";
    main [shape=box];
    main -> parse [weight=8];
    parse -> execute;
    main -> init [style=dotted];
    main -> cleanup;
    execute -> { make_string; printf};
    init -> make_string;
    edge [color=red];
    main -> printf [style=bold,label="100 times"];
    make_string [label="make a string"];
    node [shape=box,style=filled,color=".7 .3 1.0"];
    execute -> compare;
  }
custom_mark10
</details>
Fig. 1 Initial system. API based architecture

After the company was merged with the two new clients, the API service will serve to all the customer requests, total 3.
So, this increases the requests load into the API endpoint. Customers could experience some delays in the response due the API has more work to process. It maybe also some requests can be lost since the API does not implement a queue system. Fig. 2.

![Alt text](http://jkin.be/content/WallStreet-Img11.png)
<details> 
<summary></summary>
custom_mark10
  digraph G {
    size ="4,4";
    main [shape=box];
    main -> parse [weight=8];
    parse -> execute;
    main -> init [style=dotted];
    main -> cleanup;
    execute -> { make_string; printf};
    init -> make_string;
    edge [color=red];
    main -> printf [style=bold,label="100 times"];
    make_string [label="make a string"];
    node [shape=box,style=filled,color=".7 .3 1.0"];
    execute -> compare;
  }
custom_mark11
</details> 
Fig. 2 Current Situation. API based architecture

Proposed Patterns
Key issues identified:
•	Overloading to the API service.
•	Escalation limitation.
•	Delay on delivering customer’s responses.
The stock exchange applications must be robust and have a good implementation since the data handled needs real-time interaction.
So, for this problem we analyze the next Integration patterns approaches.
Messaging Patterns – Message Dispatcher

![Alt text](http://jkin.be/content/MessageDispatcher.gif)
<details> 
<summary></summary>
custom_mark10
  digraph G {
    size ="4,4";
    main [shape=box];
    main -> parse [weight=8];
    parse -> execute;
    main -> init [style=dotted];
    main -> cleanup;
    execute -> { make_string; printf};
    init -> make_string;
    edge [color=red];
    main -> printf [style=bold,label="100 times"];
    make_string [label="make a string"];
    node [shape=box,style=filled,color=".7 .3 1.0"];
    execute -> compare;
  }
custom_mark12
</details>
Fig. 3 Message Dispatcher architecture

Create a Message Dispatcher on a channel that will consume messages from a channel and distribute them to performers.
An application is using Messaging. The application needs multiple consumers on a single Message Channel to work in a coordinated fashion.

Under this approach, the Wall Street system will be able to work with several clients that consumes from a channel.

As the data that goes over the channel is not intended to be consumed by all the consumers, we need to think in a approach that can differentiate the messages.

Messaging Patterns – Selective Consumer
![Alt text](http://jkin.be/content/MessageSelectorSolution.gif)
<details> 
<summary></summary>
custom_mark10
  digraph G {
    size ="4,4";
    main [shape=box];
    main -> parse [weight=8];
    parse -> execute;
    main -> init [style=dotted];
    main -> cleanup;
    execute -> { make_string; printf};
    init -> make_string;
    edge [color=red];
    main -> printf [style=bold,label="100 times"];
    make_string [label="make a string"];
    node [shape=box,style=filled,color=".7 .3 1.0"];
    execute -> compare;
  }
custom_mark13
</details> 
Fig. 4 Selective Consumer architecture

Make the consumer a Selective Consumer, one that filteres the messages delivered by its channel so that it only receives the ones that match its criteria.
An application is using Messaging. It consumes Messages from a Message Channel, but it does not necessarily want to consume all of the messages on that channel, just some of them.

The producer will generate different types of messages and will broadcast to the channel, consumers that are subscribed for certain types of messages will receive it, and will ignore others, this way we pretend to increase the data flow and increase the response times, making the consumers can obtain data in a high rates.

