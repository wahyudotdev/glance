/**
 * Copies text to the clipboard using the modern Clipboard API if available,
 * falling back to a legacy textarea method for non-secure contexts (HTTP).
 */
export async function copyToClipboard(text: string): Promise<boolean> {
  // Try modern API first
  if (navigator.clipboard && window.isSecureContext) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch (err) {
      console.error('Clipboard API failed, falling back...', err);
    }
  }

  // Fallback for non-secure contexts (HTTP)
  try {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    
    // Ensure it's not visible and doesn't affect layout
    textArea.style.position = 'fixed';
    textArea.style.left = '-9999px';
    textArea.style.top = '0';
    textArea.style.opacity = '0';
    document.body.appendChild(textArea);
    
    textArea.focus();
    textArea.select();

    const successful = document.execCommand('copy');
    document.body.removeChild(textArea);
    return successful;
  } catch (err) {
    console.error('Fallback copy failed:', err);
    return false;
  }
}
