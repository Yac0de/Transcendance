import 'vue-i18n';
import { DefineComponent } from 'vue';

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $t: (key: string, ...args: unknown[]) => string;
    $i18n: unknown;
  }
}