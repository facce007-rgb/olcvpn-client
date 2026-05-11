package com.olc.vpn.ui.theme

import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.ui.graphics.Color

// OLC VPN Dark Theme Colors
val md_theme_dark_primary = Color(0xFF00E5FF)
val md_theme_dark_onPrimary = Color(0xFF003544)
val md_theme_dark_primaryContainer = Color(0xFF004D61)
val md_theme_dark_onPrimaryContainer = Color(0xFF97F0FF)
val md_theme_dark_secondary = Color(0xFF00C853)
val md_theme_dark_onSecondary = Color(0xFF003910)
val md_theme_dark_secondaryContainer = Color(0xFF005319)
val md_theme_dark_onSecondaryContainer = Color(0xFF6EFF8E)
val md_theme_dark_error = Color(0xFFFF1744)
val md_theme_dark_onError = Color(0xFF690005)
val md_theme_dark_errorContainer = Color(0xFF93000A)
val md_theme_dark_onErrorContainer = Color(0xFFFFDAD6)
val md_theme_dark_background = Color(0xFF0D0D0D)
val md_theme_dark_onBackground = Color(0xFFE1E2E8)
val md_theme_dark_surface = Color(0xFF1A1A1A)
val md_theme_dark_onSurface = Color(0xFFE1E2E8)
val md_theme_dark_surfaceVariant = Color(0xFF40484C)
val md_theme_dark_onSurfaceVariant = Color(0xFFBFC8CC)

val DarkColorScheme = darkColorScheme(
    primary = md_theme_dark_primary,
    onPrimary = md_theme_dark_onPrimary,
    primaryContainer = md_theme_dark_primaryContainer,
    onPrimaryContainer = md_theme_dark_onPrimaryContainer,
    secondary = md_theme_dark_secondary,
    onSecondary = md_theme_dark_onSecondary,
    secondaryContainer = md_theme_dark_secondaryContainer,
    onSecondaryContainer = md_theme_dark_onSecondaryContainer,
    error = md_theme_dark_error,
    onError = md_theme_dark_onError,
    errorContainer = md_theme_dark_errorContainer,
    onErrorContainer = md_theme_dark_onErrorContainer,
    background = md_theme_dark_background,
    onBackground = md_theme_dark_onBackground,
    surface = md_theme_dark_surface,
    onSurface = md_theme_dark_onSurface,
    surfaceVariant = md_theme_dark_surfaceVariant,
    onSurfaceVariant = md_theme_dark_onSurfaceVariant,
)

// Light theme (optional)
val LightColorScheme = lightColorScheme(
    primary = Color(0xFF006780),
    onPrimary = Color(0xFFFFFFFF),
    primaryContainer = Color(0xFFB8EAFF),
    onPrimaryContainer = Color(0xFF001F28),
    secondary = Color(0xFF006C2E),
    onSecondary = Color(0xFFFFFFFF),
    secondaryContainer = Color(0xFF8EF8A6),
    onSecondaryContainer = Color(0xFF002108),
    error = Color(0xFFBA1A1A),
    onError = Color(0xFFFFFFFF),
    errorContainer = Color(0xFFFFDAD6),
    onErrorContainer = Color(0xFF410002),
    background = Color(0xFFFBFCFE),
    onBackground = Color(0xFF191C1E),
    surface = Color(0xFFFBFCFE),
    onSurface = Color(0xFF191C1E),
    surfaceVariant = Color(0xFFDCE4E8),
    onSurfaceVariant = Color(0xFF40484C),
)
