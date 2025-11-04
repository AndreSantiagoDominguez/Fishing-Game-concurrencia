# Juego de Pesca Concurrente en Go

## Descripción del Proyecto

Aplicación interactiva de simulación de pesca desarrollada en Go que demuestra la implementación práctica de conceptos avanzados de programación concurrente. El juego permite al jugador moverse alrededor de un lago, lanzar su anzuelo y capturar peces de diferentes tipos y rareza, mientras múltiples goroutines gestionan el comportamiento independiente de cada entidad.

---

## Características Principales

El sistema implementa un entorno de pesca donde cada pez opera mediante su propia goroutine independiente, nadando de manera autónoma por el lago. Un sistema de productor-consumidor coordina la generación continua de nuevos peces y el procesamiento de capturas. Los peces se clasifican en cuatro categorías de rareza: comunes, raros, épicos y legendarios, cada uno con diferentes probabilidades de aparición y valores de puntuación.

El juego incluye un sistema de límites que mantiene el balance poblacional del lago, evitando saturación mientras preserva las probabilidades originales de aparición. Los peces tienen un tiempo de vida limitado de treinta segundos, después del cual desaparecen automáticamente, creando una dinámica de juego más interesante. El sistema proporciona retroalimentación visual mediante parpadeo durante los últimos cinco segundos antes de que un pez desaparezca.

La interfaz muestra estadísticas en tiempo real incluyendo puntuación total, número de peces capturados por tipo de rareza, y conteo actual de peces en el lago. Los controles son intuitivos con movimiento mediante teclas WASD, lanzamiento de anzuelo con espacio, y recogida con la tecla R.

---

## Requisitos del Sistema

El proyecto requiere Go versión 1.20 o superior instalado en el sistema. La librería gráfica Ebiten v2 se descarga automáticamente mediante el sistema de módulos de Go. El juego está diseñado para ejecutarse en sistemas de escritorio con soporte para ventanas gráficas, siendo compatible con Windows, macOS y Linux.

---

## Instalación y Configuración

Para comenzar con el proyecto, primero debe clonarse el repositorio en la máquina local. Una vez clonado, es necesario navegar al directorio del proyecto mediante la terminal o línea de comandos.

El siguiente paso consiste en descargar todas las dependencias necesarias. El sistema de módulos de Go se encarga automáticamente de obtener Ebiten y cualquier otra librería requerida cuando se ejecuta el comando de instalación de dependencias.
```bash
git clone [URL-del-repositorio]
cd fishing-game
go mod download
```

---

## Estructura del Proyecto

El proyecto está organizado en múltiples archivos que separan claramente las responsabilidades de cada componente. El archivo principal actúa como punto de entrada, configurando la ventana de juego y delegando el control a la estructura principal del juego.

El directorio game contiene todos los archivos relacionados con la lógica del juego. El archivo game.go define la estructura principal que gestiona el estado global, coordina las goroutines y maneja el loop de actualización y renderizado. El archivo spawner.go implementa el patrón Productor-Consumidor, conteniendo la lógica de generación de peces y procesamiento de capturas.

Los archivos player.go, fish.go y bobber.go encapsulan respectivamente el comportamiento del jugador, los peces y el anzuelo. Cada uno maneja sus propios sprites, animaciones y lógica de actualización. El directorio assets almacena todos los recursos gráficos utilizados en el juego, incluyendo los sprite sheets del pescador, los diferentes tipos de peces, el bobber y el escenario del lago.
```
fishing-game/
├── main.go
├── go.mod
├── go.sum
├── README.md
├── assets/
│   ├── fisherman_walk_up.png
│   ├── fisherman_walk_down.png
│   ├── fisherman_walk_left.png
│   ├── fisherman_walk_right.png
│   ├── fisherman_fishing.png
│   ├── fish_common.png
│   ├── fish_rare.png
│   ├── fish_epic.png
│   ├── fish_legendary.png
│   ├── bobber.png
│   └── lake_scene.png
└── game/
    ├── game.go
    ├── spawner.go
    ├── player.go
    ├── fish.go
    └── bobber.go
```

---

## Compilación y Ejecución

Para ejecutar el juego en modo desarrollo, simplemente utilice el comando run de Go desde el directorio raíz del proyecto. Esto compilará y ejecutará la aplicación en un solo paso.
```bash
go run main.go
```

Si desea compilar un ejecutable independiente, utilice el comando build. Esto generará un archivo binario que puede ejecutarse sin necesidad de tener Go instalado en la máquina destino.
```bash
go build -o fishing-game
./fishing-game
```

Para sistemas Windows, el comando de compilación sería similar pero generaría un archivo con extensión exe.
```bash
go build -o fishing-game.exe
fishing-game.exe
```

---

## Verificación de Condiciones de Carrera

Una de las características más importantes del proyecto es la ausencia de condiciones de carrera. Para verificar esto, Go proporciona una herramienta integrada que detecta automáticamente data races durante la ejecución. El flag -race instrumenta el código para rastrear todos los accesos a memoria compartida y reportar cualquier acceso concurrente inseguro.
```bash
go run -race main.go
```

La ejecución con el detector de condiciones de carrera es ligeramente más lenta debido a la instrumentación adicional, pero no debería reportar ningún problema si el código está correctamente sincronizado. Cualquier condición de carrera detectada se imprimirá en la consola con información detallada sobre las goroutines involucradas y las líneas de código problemáticas.

---

## Controles del Juego

El jugador se controla mediante el teclado con un esquema de teclas intuitivo. Las teclas WASD permiten el movimiento en las cuatro direcciones: W para mover hacia arriba, S hacia abajo, A hacia la izquierda y D hacia la derecha. Alternativamente, también pueden usarse las teclas de flecha direccionales.

Para pescar, el jugador debe posicionarse cerca de la orilla del lago. Una vez en posición, presionar la tecla Espacio lanzará el anzuelo hacia el agua. El anzuelo permanecerá activo hasta que capture un pez o el jugador decida recogerlo. Para recoger el anzuelo sin capturar nada, presione la tecla R.

Cuando el anzuelo toca un pez, se produce automáticamente la captura. El sistema mostrará brevemente una animación de captura y luego actualizará las estadísticas del jugador. Después de aproximadamente un segundo, el control regresará al jugador para continuar pescando.

---

## Mecánicas del Juego

El lago contiene peces de cuatro tipos diferentes que aparecen con probabilidades específicas. Los peces comunes son los más frecuentes con sesenta por ciento de probabilidad de aparición, otorgando diez puntos al ser capturados. Los peces raros aparecen con veinticinco por ciento de probabilidad y valen veinticinco puntos.

Los peces épicos son considerablemente más raros con doce por ciento de probabilidad de aparición, otorgando cincuenta puntos al jugador. Finalmente, los peces legendarios son extremadamente raros con solo tres por ciento de probabilidad, pero recompensan con cien puntos al ser capturados.

El sistema implementa límites poblacionales para cada tipo de pez. Pueden existir simultáneamente hasta quince peces comunes, ocho raros, cuatro épicos y dos legendarios. Estos límites previenen la saturación del lago manteniendo el juego balanceado. Cuando se alcanza el límite de un tipo específico, el spawner simplemente espera hasta que haya espacio disponible antes de generar más.

Cada pez tiene un tiempo de vida de treinta segundos. Durante los últimos cinco segundos antes de desaparecer, el pez parpadea visualmente para advertir al jugador. Esta mecánica añade presión temporal y hace que el jugador deba priorizar qué peces capturar primero, especialmente los de mayor rareza.

---

## Implementación Técnica de Concurrencia

El proyecto implementa dos patrones principales de concurrencia que trabajan en conjunto. El patrón Productor-Consumidor gestiona la generación y procesamiento de entidades, mientras que el patrón de Workers Independientes permite que cada pez opere autónomamente.

### Patrón Productor-Consumidor

La goroutine fishSpawner actúa como productor, ejecutándose continuamente en segundo plano. Cada tres segundos intenta generar un nuevo pez, primero determinando aleatoriamente el tipo basándose en las probabilidades configuradas, luego verificando si hay espacio disponible para ese tipo específico. Si se cumplen las condiciones, crea el pez y lo envía al canal spawnChan.

El método Update del juego actúa como consumidor, leyendo del canal spawnChan en cada frame mediante un select no bloqueante. Cuando recibe un pez, lo integra a la lista de entidades activas y lanza su goroutine de movimiento. Este diseño desacopla completamente la generación de la integración, permitiendo que ambos procesos operen a diferentes ritmos.

Un segundo canal catchChan maneja las capturas. Cuando el anzuelo colisiona con un pez, se envía el tipo al canal. La goroutine catchProcessor lee continuamente de este canal, calculando los puntos y actualizando las estadísticas de manera thread-safe. Esta arquitectura asíncrona evita que el procesamiento de capturas bloquee el loop principal del juego.

### Workers Independientes

Cada pez es controlado por su propia goroutine ejecutando el método Swim. Esta goroutine mantiene un loop con ticker que actualiza la posición aproximadamente sesenta veces por segundo. El pez cambia aleatoriamente de dirección cada dos segundos con treinta por ciento de probabilidad, y rebota automáticamente cuando se acerca a los bordes del lago.

La goroutine también gestiona el tiempo de vida del pez, verificando constantemente cuánto tiempo ha transcurrido desde su creación. Cuando alcanza los treinta segundos, marca su estado como inactivo y termina su ejecución limpiamente. Esta auto-gestión elimina la necesidad de lógica externa para limpieza de entidades.

### Sincronización y Protección de Datos

Todo el acceso a datos compartidos está protegido por mutex. La estructura Game tiene un mutex que protege la lista de peces activos, estadísticas de puntuación y otros datos globales. Cada pez también tiene su propio mutex protegiendo su posición y estado interno, permitiendo que múltiples goroutines de peces operen simultáneamente sin interferencia.

Las operaciones con canales se realizan mediante select con caso default, haciéndolas no bloqueantes. Esto es crucial en Update y Draw donde no se pueden permitir bloqueos que comprometan la fluidez visual. El uso de defer para liberar mutex asegura que siempre se liberen incluso ante errores inesperados.

### Gestión del Ciclo de Vida

Un Context creado con context.WithCancel propaga señales de cancelación a todas las goroutines. Cada goroutine de larga duración incluye un case en su select que escucha el canal Done del context, terminando limpiamente cuando recibe la señal de cierre.

Un WaitGroup rastrea todas las goroutines activas. Antes de lanzar cualquier goroutine, se incrementa el contador con Add, y la goroutine llama a Done mediante defer al terminar. El método Cleanup cancela el context, cierra los canales y espera en el WaitGroup, asegurando un cierre ordenado sin fugas de recursos.

---

## Arquitectura de Renderizado

La separación entre lógica y renderizado es fundamental para mantener la fluidez visual. El método Update ejecuta toda la lógica de actualización del juego, incluyendo lectura de canales, actualización de entidades y detección de colisiones. Todas estas operaciones son rápidas y deterministas, completándose en fracción de frame.

El método Draw es exclusivamente para renderizado. Adquiere el mutex brevemente solo para leer las referencias necesarias, luego lo libera inmediatamente antes de dibujar. Esto minimiza el tiempo de bloqueo, permitiendo que otras goroutines continúen trabajando durante el renderizado. Draw nunca modifica el estado del juego ni ejecuta operaciones bloqueantes.

Los sprites se cargan una vez al inicio y se almacenan en estructuras globales compartidas. Todas las instancias de peces usan las mismas imágenes base, aplicando transformaciones en tiempo de renderizado según su estado. Este patrón flyweight reduce significativamente el uso de memoria cuando hay múltiples entidades similares.

---

## Prevención de Errores Comunes

El proyecto implementa varias estrategias para prevenir problemas típicos de programación concurrente. La más importante es la prevención de capturas múltiples. Cuando se detecta una colisión, el bobber se desactiva inmediatamente estableciendo su campo active en false antes de iniciar cualquier procesamiento adicional. Esto previene que frames subsecuentes detecten colisiones adicionales antes del reset.

La verificación de límites de peces se realiza dentro del mutex del juego. Sin esta protección, invocaciones concurrentes del spawner podrían todas leer el mismo conteo y decidir generar peces simultáneamente, excediendo el límite. El mutex serializa estas verificaciones garantizando consistencia.

Todas las goroutines tienen condiciones de salida claras. Ya sea por señal del context, por expiración de tiempo de vida, o por desactivación explícita, cada goroutine puede terminar limpiamente sin quedarse bloqueada indefinidamente. Los canales se cierran solo después de cancelar el context, permitiendo que goroutinas bloqueadas se desbloqueen y verifiquen si deben terminar.

---

## Estadísticas y Métricas

La interfaz muestra información completa sobre el estado del juego. La puntuación total acumulada se actualiza en tiempo real con cada captura. Los contadores por tipo de pez muestran cuántos de cada rareza ha capturado el jugador durante la sesión.

Una línea adicional muestra cuántos peces de cada tipo están actualmente nadando en el lago. Esta información permite al jugador observar visualmente el funcionamiento del sistema de límites y la distribución de probabilidades. Los sufijos C, R, E y L representan respectivamente comunes, raros, épicos y legendarios.

Toda esta información se actualiza de manera thread-safe. Las goroutines que modifican estadísticas adquieren el mutex antes de hacer cambios, y el método Draw adquiere el mutex brevemente para leer los valores actuales antes de mostrarlos. Esto garantiza que nunca se muestren valores inconsistentes o corruptos.

---

## Solución de Problemas

Si el juego no encuentra los archivos de sprites, verifique que el directorio assets existe en la ubicación correcta relativa al ejecutable. Todos los archivos PNG deben estar presentes con los nombres exactos especificados en el código. El juego puede funcionar sin el escenario del lago mostrando un fondo de color sólido, pero requiere los sprites de peces y del jugador.

Si experimenta problemas de rendimiento, verifique que no esté ejecutando el juego con el flag -race a menos que específicamente esté probando condiciones de carrera. El detector de races añade overhead significativo que puede afectar la fluidez en sistemas de recursos limitados. La compilación normal sin instrumentación debería ejecutarse fluidamente en hardware moderno.

En caso de que el juego se cierre inesperadamente, ejecute con el flag -race para verificar si existe alguna condición de carrera no detectada. Revise también la consola para mensajes de error o panics que puedan indicar el problema. Los errores de carga de sprites se reportan claramente en la inicialización.

---
# Fishing-Game-concurrencia
