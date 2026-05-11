package com.olc.vpn.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.Stop
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
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

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        // Status Card
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(
                containerColor = when (connectionState) {
                    "connected" -> MaterialTheme.colorScheme.primaryContainer
                    "connecting" -> MaterialTheme.colorScheme.secondaryContainer
                    "error" -> MaterialTheme.colorScheme.errorContainer
                    else -> MaterialTheme.colorScheme.surfaceVariant
                }
            )
        ) {
            Column(
                modifier = Modifier.padding(16.dp),
                horizontalAlignment = Alignment.CenterHorizontally
            ) {
                Text(
                    text = connectionState.uppercase(),
                    style = MaterialTheme.typography.headlineMedium
                )
                Spacer(modifier = Modifier.height(8.dp))
                Text(
                    text = selectedProfile ?: "No profile selected",
                    style = MaterialTheme.typography.bodyMedium
                )
            }
        }

        // Connect/Disconnect Button
        Button(
            onClick = {
                if (connectionState == "connected") {
                    viewModel.disconnect()
                } else {
                    onConnectClick()
                }
            },
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
            enabled = connectionState != "connecting"
        ) {
            Icon(
                imageVector = if (connectionState == "connected") Icons.Default.Stop else Icons.Default.PlayArrow,
                contentDescription = null
            )
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                text = if (connectionState == "connected") "Disconnect" else "Connect",
                style = MaterialTheme.typography.titleMedium
            )
        }

        // Traffic Stats
        Card(
            modifier = Modifier.fillMaxWidth()
        ) {
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Text(
                    text = "Traffic Statistics",
                    style = MaterialTheme.typography.titleMedium
                )
                Spacer(modifier = Modifier.height(8.dp))
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Column {
                        Text("Upload", style = MaterialTheme.typography.bodySmall)
                        Text(
                            formatBytes(bytesUp),
                            style = MaterialTheme.typography.titleLarge
                        )
                    }
                    Column(horizontalAlignment = Alignment.End) {
                        Text("Download", style = MaterialTheme.typography.bodySmall)
                        Text(
                            formatBytes(bytesDown),
                            style = MaterialTheme.typography.titleLarge
                        )
                    }
                }
            }
        }

        // Quick Actions
        Card(
            modifier = Modifier.fillMaxWidth()
        ) {
            Column(
                modifier = Modifier.padding(16.dp)
            ) {
                Text(
                    text = "Quick Actions",
                    style = MaterialTheme.typography.titleMedium
                )
                Spacer(modifier = Modifier.height(8.dp))
                OutlinedButton(
                    onClick = { /* TODO: Select profile */ },
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Text("Select Profile")
                }
                Spacer(modifier = Modifier.height(8.dp))
                OutlinedButton(
                    onClick = { /* TODO: Import profile */ },
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Text("Import Profile")
                }
            }
        }
    }
}

private fun formatBytes(bytes: Long): String {
    return when {
        bytes < 1024 -> "${bytes}B"
        bytes < 1024 * 1024 -> String.format("%.1fKB", bytes / 1024.0)
        bytes < 1024 * 1024 * 1024 -> String.format("%.1fMB", bytes / (1024.0 * 1024))
        else -> String.format("%.1fGB", bytes / (1024.0 * 1024 * 1024))
    }
}
