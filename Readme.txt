Cet executable effectue les opérations suivantes :
1- Ouvre les images .jpg et .png dans le dosser "in"
2- Les convertit en .png et les place dans le dossier "out"
3- Les réduit avec une largeur max de 320 px et les place dans "resized"
4- Optimize leur taille et les place dans "compressed"

Lancez simplement "imgformatter.exe" pour lancer les traitements

Pour spécifier des paramètres :
dans une invite de commande, lancer 
imgformatter.exe -width=400 -height=9999 -cropwidth=200 -cropheight=150

paramètres :
width : largeur minimum de l'image (défaut: 320)
height : hauteur minimum de l'image (défaut: 500)

cropwidth et cropheight : crop (tronque) l'image pour qu'elle aie des dimensions précises, par exemple 320 * 320 px
si un des deux paramètres n'est pas précisé, l'image ne sera pas tronquée