package com.olc.vpn.service

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Intent
import android.net.VpnService
import android.os.Build
import android.os.ParcelFileDescriptor
import androidx.core.app.NotificationCompat
import com.olc.vpn.MainActivity
import kotlinx.coroutines.*
import mobile.Mobile
import mobile.VPNCore

class OlcVpnService : VpnService() {

    private var vpnInterface: ParcelFileDescriptor? = null
    private var vpnCore: VPNCore? = null
    private val serviceScope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    companion object {
        private const val NOTIFICATION_ID = 1
        private const val CHANNEL_ID = "olc_vpn_channel"
        const val ACTION_CONNECT = "com.olc.vpn.CONNECT"
        const val ACTION_DISCONNECT = "com.olc.vpn.DISCONNECT"
        const val EXTRA_PROFILE_JSON = "profile_json"
    }

    override fun onCreate() {
        super.onCreate()
        createNotificationChannel()

        vpnCore = Mobile.newVPNCore()
        val dataDir = applicationContext.filesDir.absolutePath
        vpnCore?.initialize(dataDir)
    }

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        when (intent?.action) {
            ACTION_CONNECT -> {
                val profileJson = intent.getStringExtra(EXTRA_PROFILE_JSON)
                if (profileJson != null) {
                    startForeground(NOTIFICATION_ID, createNotification("Connecting..."))
                    connect(profileJson)
                }
            }

            ACTION_DISCONNECT -> {
                disconnect()
                stopSelf()
            }
        }

        return START_STICKY
    }

    private fun connect(profileJson: String) {
        serviceScope.launch {
            try {
                val builder = Builder()
                    .setSession("OLC VPN")
                    .addAddress("10.0.0.1", 30)
                    .addRoute("0.0.0.0", 0)
                    .addDnsServer("1.1.1.1")
                    .addDnsServer("8.8.8.8")
                    .setMtu(1500)
                    .setBlocking(false)

                vpnInterface = builder.establish()

                if (vpnInterface == null) {
                    updateNotification("Connection failed")
                    stopSelf()
                    return@launch
                }

                val fd = vpnInterface!!.fd

                vpnCore?.startWithTunFd(fd.toLong())

                updateNotification("Connected")

                monitorTraffic()

            } catch (e: Exception) {
                e.printStackTrace()
                updateNotification("Error: ${e.message}")
                disconnect()
            }
        }
    }

    private fun disconnect() {
        serviceScope.launch {
            try {
                vpnCore?.disconnect()
                vpnInterface?.close()
                vpnInterface = null
                updateNotification("Disconnected")
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
    }

    private suspend fun monitorTraffic() {
        while (vpnInterface != null) {
            delay(1000)

            val bytesUp = vpnCore?.bytesUp ?: 0
            val bytesDown = vpnCore?.bytesDown ?: 0
            val status = vpnCore?.status ?: "unknown"

            if (status == "connected") {
                updateNotification(
                    "Connected - ↑${formatBytes(bytesUp)} ↓${formatBytes(bytesDown)}"
                )
            }
        }
    }

    private fun formatBytes(bytes: Long): String {
        return when {
            bytes < 1024 -> "${bytes}B"
            bytes < 1024 * 1024 -> "${bytes / 1024}KB"
            bytes < 1024 * 1024 * 1024 -> "${bytes / (1024 * 1024)}MB"
            else -> "${bytes / (1024 * 1024 * 1024)}GB"
        }
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                CHANNEL_ID,
                "VPN Service",
                NotificationManager.IMPORTANCE_LOW
            ).apply {
                description = "OLC VPN connection status"
                setShowBadge(false)
            }

            val notificationManager = getSystemService(NotificationManager::class.java)
            notificationManager.createNotificationChannel(channel)
        }
    }

    private fun createNotification(text: String): Notification {
        val intent = Intent(this, MainActivity::class.java)

        val pendingIntent = PendingIntent.getActivity(
            this,
            0,
            intent,
            PendingIntent.FLAG_IMMUTABLE or PendingIntent.FLAG_UPDATE_CURRENT
        )

        return NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle("OLC VPN")
            .setContentText(text)
            .setSmallIcon(android.R.drawable.stat_sys_download_done)
            .setContentIntent(pendingIntent)
            .setOngoing(true)
            .build()
    }

    private fun updateNotification(text: String) {
        val notificationManager = getSystemService(NotificationManager::class.java)
        notificationManager.notify(
            NOTIFICATION_ID,
            createNotification(text)
        )
    }

    override fun onDestroy() {
        super.onDestroy()
        disconnect()
        serviceScope.cancel()
    }
}
