Cet executable effectue les opérations suivantes :
1- Ouvre les images .jpg et .png dans le dosser "in"
2- Les convertit en .png et les places dans le dossier "out"
3- Les réduit avec une largeur max de 320 px et les place dans "resized"
4- Optimize leur taille et les place dans "compressed"

Lancez simplement "imgformatter.exe" pour lancer les traitements

Pour spécifier des paramètres :
dans une invite de commande, lancer 
imgformatter.exe -width=320 -height=500 -cropwidth=320 -cropheight=320

paramètres :
width : largeur minimum de l'image (0 = infini)
height : hauteur minimum de l'image (0 = infini)

cropwidth et cropheight : crop (tronque) l'image pour qu'elle aie des dimensions précises, par exemple 320 * 320 px
si un des deux paramètres n'est pas précisé, l'image ne sera pas tronquée