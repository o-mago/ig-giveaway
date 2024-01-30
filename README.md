![Captura de tela de 2024-01-30 16-55-59](https://github.com/o-mago/ig-giveaway/assets/23153316/3f480193-1852-4fc5-86a1-44a27f8fe205)

## About

`ig-giveaway` is a cli tool to do giveaways based on comments in an Instagram post

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

To use the tool, you must create a meta developer account and an app with the Instagram's Graph API access: `https://developers.facebook.com/apps/creation`

Using the [Graph Api Explorer](https://developers.facebook.com/tools/explorer), you must add these scopes: `pages_show_list,instagram_basic,pages_read_engagement` and click `Generate Access Token`
Use this token to get

![graph](https://github.com/o-mago/ig-giveaway/assets/23153316/3d107704-c36b-4fa4-a03d-0d166cbd3b7b)


