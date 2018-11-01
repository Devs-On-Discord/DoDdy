package com.github.dod.doddy.guilds

data class Guild(val id: String,
                 val botPrefix: Char = '!',
                 val welcomeChannel: String,
                 val announcementChannel: String,
                 val votingChannel: String,
                 val roleGradient1: Color,
                 val roleGradient2: Color)