package com.olc.vpn.ui.screens

import androidx.compose.foundation.Canvas
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.olc.vpn.viewmodel.VpnViewModel

// HomeScreen в стиле Hiddify/v2RayTun
@Composable
fun HomeScreen(
    viewModel: VpnViewModel,
    onConnectClick: () -> Unit
) {
    val connectionState by viewModel.connectionState.collectAsState()
    val bytesUp by viewModel.bytesUp.collectAsState()
    val bytesDown by viewModel.bytesDown.collectAsState()
    val selectedProfile by viewModel.selectedProfile.collectAsState()
    val latency by viewModel.latency.collectAsState()

    // Цвета как в Hiddify
    val statusColor = when (connectionState) {
        "connected" -> Color(0xFF2E7D32)      // Green 800
        "connecting" -> Color(0xFFFFC107)     // Amber
        "error" -> Color(0xFFD32F2F)          // Red 700
        else -> Color(0xFF3F51B5)             // Indigo (disconnected)
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        // Карточка профиля сверху (как в Hiddify)
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(bottom = 16.dp),
            colors = CardDefaults.cardColors(
                containerColor = MaterialTheme.colorScheme.surfaceVariant
            )
        ) {
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Text(
                    text = selectedProfile ?: "No Profile Selected",
                    style = MaterialTheme.typography.titleMedium,
                    fontWeight = FontWeight.Bold
                )
                if (selectedProfile != null) {
                    Text(
                        text = "VLESS • Reality",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurfaceVariant
                    )
                }
            }
        }

        Spacer(modifier = Modifier.weight(1f))

        // Большая круглая кнопка подключения (как в Hiddify)
        Box(
            contentAlignment = Alignment.Center,
            modifier = Modifier.size(200.dp)
        ) {
            // Круг с обводкой
            Canvas(modifier = Modifier.fillMaxSize()) {
                drawCircle(
                    color = statusColor,
                    style = Stroke(width = 16.dp.toPx())
                )
                if (connectionState == "connected" || connectionState == "connecting") {
                    drawCircle(
                        color = statusColor.copy(alpha = 0.3f)
                    )
                }
            }

            // Текст статуса в центре
            Column(
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = when (connectionState) {
                        "connected" -> "Connected"
                        "connecting" -> "Connecting..."
                        "error" -> "Failed"
                        else -> "Tap to Connect"
                    },
                    style = MaterialTheme.typography.titleLarge,
                    fontWeight = FontWeight.Bold
                )
            }

            // Кликабельная область
            Surface(
                onClick = {
                    if (connectionState == "connected") {
                        viewModel.disconnect()
                    } else if (connectionState != "connecting") {
                        onConnectClick()
                    }
                },
                modifier = Modifier.fillMaxSize(),
                color = Color.Transparent,
                shape = CircleShape
            ) {}
        }

        // Индикатор задержки под кругом (как в Hiddify)
        Spacer(modifier = Modifier.height(24.dp))
        Divider()
        Spacer(modifier = Modifier.height(8.dp))
        Text(
            text = if (latency > 0 && latency < 65000) "${latency} ms" else "- ms",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = if (latency > 0 && latency < 200) Color(0xFF2E7D32) else MaterialTheme.colorScheme.onSurface
        )

        Spacer(modifier = Modifier.weight(1f))

        // Футер с метриками (как в Hiddify)
        Divider()
        Spacer(modifier = Modifier.height(8.dp))
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.Center,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = "↑ ${formatBytes(bytesUp)}",
                style = MaterialTheme.typography.bodyMedium
            )
            Text(
                text = "  •  ",
                style = MaterialTheme.typography.bodyMedium
            )
            Text(
                text = "↓ ${formatBytes(bytesDown)}",
                style = MaterialTheme.typography.bodyMedium
            )
        }
        Spacer(modifier = Modifier.height(16.dp))
    }
}

private fun formatBytes(bytes: Long): String {
    return when {
        bytes < 1024 -> "${bytes} B"
        bytes < 1024 * 1024 -> String.format("%.1f KB", bytes / 1024.0)
        bytes < 1024 * 1024 * 1024 -> String.format("%.1f MB", bytes / (1024.0 * 1024))
        else -> String.format("%.1f GB", bytes / (1024.0 * 1024 * 1024))
    }
}
