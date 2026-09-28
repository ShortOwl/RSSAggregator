export const themeKey = "margin.theme";
export type Theme = "light" | "dark" | "system";

export function parseTheme(value: string | null): Theme {
  return value === "light" || value === "dark" ? value : "system";
}

// Runs in the document head before any content can paint.
export const themeScript = `(function(){var t="system";try{var v=localStorage.getItem("${themeKey}");if(v==="light"||v==="dark")t=v}catch(e){}document.documentElement.dataset.theme=t==="system"?(matchMedia("(prefers-color-scheme: dark)").matches?"dark":"light"):t})()`;
