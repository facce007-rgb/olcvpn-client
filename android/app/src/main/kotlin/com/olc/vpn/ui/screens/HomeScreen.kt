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
import com.olc.vpn.viewmodel.VpnViewModel

@Composable
fun HomeScreen(
    viewModel: VpnViewModel,
    onConnectClick: () -> Unit
) {
    val connectionState by viewModel.connectionState.collectAsState()
    val bytesUp by viewModel.bytesUp.collectAsState()
    val bytesDown by viewModel.bytesDown.collectAsState()
    val selectedProfile by viewModel.selectedProfile.collectAsState()

    val latency = 0

    val statusColor = when (connectionState) {
        "connected" -> Color(0xFF2E7D32)
        "connecting" -> Color(0xFFFFC107)
        "error" -> Color(0xFFD32F2F)
        else -> Color(0xFF3F51B5)
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
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

        Box(
            contentAlignment = Alignment.Center,
            modifier = Modifier.size(200.dp)
        ) {
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

        Spacer(modifier = Modifier.height(24.dp))
        Divider()
        Spacer(modifier = Modifier.height(8.dp))
        Text(
            text = if (latency > 0 && latency < 65000) "${latency} ms" else "- ms",
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.onSurface
        )

        Spacer(modifier = Modifier.weight(1f))

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
