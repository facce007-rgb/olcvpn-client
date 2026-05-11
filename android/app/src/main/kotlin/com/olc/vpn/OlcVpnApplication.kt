package com.olc.vpn

import android.app.Application
import android.content.Context

class OlcVpnApplication : Application() {

    companion object {
        private lateinit var instance: OlcVpnApplication

        fun getContext(): Context = instance.applicationContext
    }

    override fun onCreate() {
        super.onCreate()
        instance = this
    }
}
