create table articles
(
    id int unsigned auto_increment,
    title varchar(255) not null,
    source varchar(255) not null,
    published_at date not null,
    excerpt text not null,
    thumbnail_url text not null,
    url text not null,
    primary key (title(64), published_at),
    index articles_id_index (id)
);

create table event_categories
(
    id int unsigned auto_increment
        primary key,
    name varchar(191) not null,
    pattern varchar(191) not null,
    info_url text null,
    image_url text null,
    description text null,
    priority int unsigned not null,
    constraint event_categories_name_uindex
        unique (name),
    constraint event_categories_pattern_uindex
        unique (pattern)
);

create table event_category_articles
(
    event_category_id int unsigned not null,
    article_id int unsigned not null,
    primary key (event_category_id, article_id),
    constraint event_category_articles_article_id_fk
        foreign key (article_id) references articles (id),
    constraint event_category_articles_event_category_id_fk
        foreign key (event_category_id) references event_categories (id)
);

create table events
(
    id int unsigned auto_increment
        primary key,
    event_category_id int unsigned null,
    team_liquid_id int unsigned null,
    title varchar(255) not null,
    stage varchar(255) not null,
    starts_at datetime not null,
    constraint events_team_liquid_id_uindex
        unique (team_liquid_id),
    constraint events_event_categories_id_fk
        foreign key (event_category_id) references event_categories (id)
);

create table platforms
(
    id int unsigned auto_increment
        primary key,
    name varchar(255) not null
);

create table channels
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

create table event_category_channels
(
    event_category_id int unsigned not null,
    channel_id varchar(64) not null,
    primary key (event_category_id, channel_id),
    constraint event_category_channels_channels_id_fk
        foreign key (channel_id) references channels (id),
    constraint event_category_channels_event_categories_id_fk
        foreign key (event_category_id) references event_categories (id)
);

create table players
(
    id int unsigned auto_increment
        primary key,
    player_id varchar(191) not null,
    name varchar(255) default '' not null,
    race varchar(8) not null,
    team varchar(255) default '' not null,
    country varchar(255) default '' not null,
    total_earnings decimal default 0 not null,
    date_of_birth date null,
    liquipedia_url text null,
    image_url text null,
    stream_url text null,
    is_retired tinyint(1) default 0 not null,
    constraint players_player_id_name_uindex
        unique (player_id, name(64))
);

create table player_articles
(
    player_id int unsigned not null,
    article_id int unsigned not null,
    primary key (player_id, article_id),
    constraint player_articles_article_id_fk
        foreign key (article_id) references articles (id),
    constraint player_articles_player_id_fk
        foreign key (player_id) references players (id)
);

create table videos
(
    id varchar(64) not null
        primary key,
    event_category_id int unsigned null,
    platform_id int unsigned not null,
    channel_id varchar(64) null,
    title varchar(255) not null,
    thumbnail_url text default '' not null,
    view_count int unsigned default 0 not null,
    duration int unsigned default 0 not null,
    created_at datetime default current_timestamp() not null,
    updated_at datetime default current_timestamp() not null,
    type varchar(64) default '' null,
    constraint videos_channels_id_fk
        foreign key (channel_id) references sc2hub.channels (id)
            on delete set null,
    constraint videos_event_categories_id_fk
        foreign key (event_category_id) references sc2hub.event_categories (id)
);

create table player_videos
(
    player_id int unsigned not null,
    video_id varchar(64) not null,
    primary key (player_id, video_id),
    constraint player_videos_players_id_fk
        foreign key (player_id) references players (id),
    constraint player_videos_videos_id_fk
        foreign key (video_id) references videos (id)
);

create table users
(
    id int unsigned auto_increment
        primary key,
    username varchar(191) not null,
    password_hash varchar(255) not null,
    constraint users_username_uindex
        unique (username)
);
