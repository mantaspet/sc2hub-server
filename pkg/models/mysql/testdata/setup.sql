use test_sc2hub;

# create or replace table articles
# (
#     id int unsigned auto_increment,
#     title varchar(255) not null,
#     source varchar(255) not null,
#     published_at date not null,
#     excerpt text not null,
#     thumbnail_url text not null,
#     url text not null,
#     primary key (title, published_at)
# );
#
# create or replace index if not exists articles_id_index
#     on articles (id);

create or replace table event_categories
(
    id int unsigned auto_increment
        primary key,
    name varchar(255) not null,
    pattern varchar(255) not null,
    info_url text default '' not null,
    image_url text default '' not null,
    description text default '' not null,
    priority int unsigned not null,
    constraint event_categories_name_uindex
        unique (name),
    constraint event_categories_pattern_uindex
        unique (pattern)
);

# create or replace table event_category_articles
# (
#     event_category_id int unsigned not null,
#     article_id int unsigned not null,
#     primary key (event_category_id, article_id),
#     constraint event_category_articles_article_id_fk
#         foreign key (article_id) references articles (id),
#     constraint event_category_articles_event_category_id_fk
#         foreign key (event_category_id) references event_categories (id)
# );

# create or replace table events
# (
#     id int unsigned auto_increment
#         primary key,
#     event_category_id int unsigned null,
#     team_liquid_id int unsigned null,
#     title varchar(255) not null,
#     stage varchar(255) not null,
#     starts_at datetime not null,
#     constraint events_team_liquid_id_uindex
#         unique (team_liquid_id),
#     constraint events_event_categories_id_fk
#         foreign key (event_category_id) references event_categories (id)
# );

create or replace table platforms
(
    id int unsigned auto_increment
        primary key,
    name varchar(255) not null
);

create or replace table channels
(
    id varchar(64) not null
        primary key,
    platform_id int unsigned not null,
    login varchar(255) default '' not null,
    title varchar(255) not null,
    profile_image_url text not null,
    constraint channels_platforms_id_fk
        foreign key (platform_id) references platforms (id)
);

create or replace table event_category_channels
(
    id int unsigned auto_increment
        primary key,
    event_category_id int unsigned not null,
    channel_id varchar(64) not null,
    constraint event_category_channels_channels_id_fk
        foreign key (channel_id) references channels (id),
    constraint event_category_channels_event_categories_id_fk
        foreign key (event_category_id) references event_categories (id)
);

# create or replace table players
# (
#     id int unsigned auto_increment
#         primary key,
#     player_id varchar(255) not null,
#     name varchar(255) default '' not null,
#     race varchar(8) not null,
#     team varchar(255) default '' not null,
#     country varchar(255) default '' not null,
#     total_earnings decimal default 0 not null,
#     date_of_birth date null,
#     liquipedia_url text default '' not null,
#     image_url text default '' not null,
#     stream_url text default '' not null,
#     is_retired tinyint(1) default 0 not null,
#     constraint players_player_id_uindex
#         unique (player_id)
# );

# create or replace table player_articles
# (
#     player_id int unsigned not null,
#     article_id int unsigned not null,
#     primary key (player_id, article_id),
#     constraint player_articles_article_id_fk
#         foreign key (article_id) references articles (id),
#     constraint player_articles_player_id_fk
#         foreign key (player_id) references players (id)
# );

# create or replace table videos
# (
#     id varchar(64) not null
#         primary key,
#     event_category_id int unsigned null,
#     platform_id int unsigned not null,
#     channel_id varchar(64) null,
#     title varchar(255) not null,
#     duration varchar(16) default '' not null,
#     thumbnail_url text default '' not null,
#     created_at datetime default current_timestamp() not null,
#     type varchar(64) default '' null,
#     constraint videos_channels_id_fk
#         foreign key (channel_id) references channels (id)
#             on delete set null,
#     constraint videos_event_categories_id_fk
#         foreign key (event_category_id) references event_categories (id)
# );

# create table player_videos
# (
#     player_id int unsigned not null,
#     video_id varchar(64) not null,
#     primary key (player_id, video_id),
#     constraint player_videos_players_id_fk
#         foreign key (player_id) references players (id),
#     constraint player_videos_videos_id_fk
#         foreign key (video_id) references videos (id)
# );

INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (1, 'World Championship Series', 'wcs', 'https://liquipedia.net/starcraft2/World_Championship_Series', 'https://static-wcs.starcraft2.com/media/images/logo/logo-event-circuit.png', '', 4);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (2, 'World Electronic Sports Games', 'wesg', '', 'https://proxy.duckduckgo.com/iu/?u=https%3A%2F%2Fstatic.hltv.org%2Fimages%2FeventLogos%2F2376.png&f=1', '', 6);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (3, 'Intel Extreme Masters', 'iem', '', 'https://proxy.duckduckgo.com/iu/?u=http%3A%2F%2Fwww.eclypsia.com%2Fpublic%2Fupload%2Fcke%2FLoL%2FIEM_Challenger_Logo.png&f=1', '', 3);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (5, 'Proxy Tempest', 'tempest', '', 'https://s3.amazonaws.com/challonge_app/organizations/images/000/061/865/hdpi/Nerazim-Tempest.png?1534435945', '', 5);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (6, 'Global Starcraft League Code S', 'gsl', '', 'https://proxy.duckduckgo.com/iu/?u=https%3A%2F%2Fliquipedia.net%2Fcommons%2Fimages%2F4%2F4a%2FGsl.png&f=1', '', 1);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (7, 'Go4SC2', 'go4sc2', '', 'http://misc.team-aaa.com/perso_Xan/go4.png', '', 2);
INSERT INTO event_categories (id, name, pattern, info_url, image_url, description, priority) VALUES (8, 'TESPA Collegiate Series', 'tespa', '', 'https://i2.wp.com/cesn.gg/wp-content/uploads/2019/02/Dz3dw0xWwAIoJoj.jpg?resize=370%2C247&ssl=1', '', 7);

INSERT INTO platforms (id, name) VALUES (1, 'twitch');
INSERT INTO platforms (id, name) VALUES (2, 'youtube');

INSERT INTO channels (id, platform_id, login, title, profile_image_url) VALUES ('42508152', 1, 'starcraft', 'StarCraft', 'https://static-cdn.jtvnw.net/jtv_user_pictures/0c9813cae3797d96-profile_image-300x300.png');
INSERT INTO channels (id, platform_id, login, title, profile_image_url) VALUES ('52229024', 1, 'gsl', 'GSL', 'https://static-cdn.jtvnw.net/jtv_user_pictures/gsl-profile_image-0670464b721ea8c5-300x300.png');
INSERT INTO channels (id, platform_id, login, title, profile_image_url) VALUES ('UCK5eBtuoj_HkdXKHNmBLAXg', 2, '', 'AfreecaTV eSports', 'https://yt3.ggpht.com/a-/AAuE7mBZ1no98oeHv-OkWsyXSL7I9Fuj9LjPZ2JcHg=s88-mo-c-c0xffffffff-rj-k-no');

INSERT INTO event_category_channels (id, event_category_id, channel_id) VALUES (1, 6, '52229024');
INSERT INTO event_category_channels (id, event_category_id, channel_id) VALUES (3, 1, '42508152');
INSERT INTO event_category_channels (id, event_category_id, channel_id) VALUES (4, 8, '42508152');
INSERT INTO event_category_channels (id, event_category_id, channel_id) VALUES (18, 6, 'UCK5eBtuoj_HkdXKHNmBLAXg');