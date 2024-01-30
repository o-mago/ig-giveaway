## About

`ig-giveaway` is a cli tool to do giveaways based on comments in an Instragram post

## Install

`go install github.com/o-mago/ig-giveaway@latest`

or

Download one of the assets from the [releases](https://github.com/o-mago/ig-giveaway/releases)

or

Build from source

## Usage

`ig-giveaway`

The options are:

- `Instagram user name`: username to whom the post belongs to
- `Instagram post code`: code from the post (can be retrieved in the post web URL)
- `Graph API Token`: access token generated using the [Graph Api Explorer](https://developers.facebook.com/tools/explorer)
- `Number of mentions`: minimum number of mentions the user must comment to be a contender
- `Should filter one entry per user?`: Each valid comment will be one entry or each user should have only one entry

## Requirements

In order to use the tool, you must create a meta developer account and a app with the instagram's Graph API access.

Using the [Graph Api Explorer](https://developers.facebook.com/tools/explorer), you must add these scopes: `pages_show_list,instagram_basic,pages_read_engagement` and click `Generate Access Token`
