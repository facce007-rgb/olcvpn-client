package com.olc.vpn.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch

class VpnViewModel : ViewModel() {

    private val _connectionState = MutableStateFlow("disconnected")
    val connectionState: StateFlow<String> = _connectionState.asStateFlow()

    private val _bytesUp = MutableStateFlow(0L)
    val bytesUp: StateFlow<Long> = _bytesUp.asStateFlow()

    private val _bytesDown = MutableStateFlow(0L)
    val bytesDown: StateFlow<Long> = _bytesDown.asStateFlow()

    private val _selectedProfile = MutableStateFlow<String?>(null)
    val selectedProfile: StateFlow<String?> = _selectedProfile.asStateFlow()

    private val _profiles = MutableStateFlow<List<String>>(emptyList())
    val profiles: StateFlow<List<String>> = _profiles.asStateFlow()

    private val _logs = MutableStateFlow<List<String>>(emptyList())
    val logs: StateFlow<List<String>> = _logs.asStateFlow()

    init {
        loadProfiles()
    }

    private fun loadProfiles() {
        viewModelScope.launch {
            // TODO: Load profiles from VPNCore
            _profiles.value = listOf(
                "Example Server 1",
                "Example Server 2"
            )
        }
    }

    fun selectProfile(profile: String) {
        _selectedProfile.value = profile
        addLog("Selected profile: $profile")
    }

    fun deleteProfile(profile: String) {
        viewModelScope.launch {
            _profiles.value = _profiles.value.filter { it != profile }
            if (_selectedProfile.value == profile) {
                _selectedProfile.value = null
            }
            addLog("Deleted profile: $profile")
        }
    }

    fun connect() {
        viewModelScope.launch {
            _connectionState.value = "connecting"
            addLog("Connecting...")
            // TODO: Start VPN service
        }
    }

    fun disconnect() {
        viewModelScope.launch {
            _connectionState.value = "disconnected"
            addLog("Disconnected")
            _bytesUp.value = 0
            _bytesDown.value = 0
            // TODO: Stop VPN service
        }
    }

    fun updateTraffic(up: Long, down: Long) {
        _bytesUp.value = up
        _bytesDown.value = down
    }

    fun updateConnectionState(state: String) {
        _connectionState.value = state
        addLog("Connection state: $state")
    }

    fun addLog(message: String) {
        val timestamp = java.text.SimpleDateFormat("HH:mm:ss", java.util.Locale.getDefault())
            .format(java.util.Date())
        _logs.value = _logs.value + "[$timestamp] $message"

        // Keep only last 100 logs
        if (_logs.value.size > 100) {
            _logs.value = _logs.value.takeLast(100)
        }
    }

    fun clearLogs() {
        _logs.value = emptyList()
    }
}
