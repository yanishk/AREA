import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:math';

class ServicesCard extends StatelessWidget {
  const ServicesCard({super.key});

  @override
  Widget build(BuildContext context) {
    List<Color> colorList = [
      const Color.fromARGB(255, 41, 44, 45),
      const Color.fromARGB(255, 255, 0, 0),
      const Color.fromARGB(255, 29, 185, 84),
    ];

    List<String> nameList = [
      "Github",
      "Youtube",
      "Spotify",
    ];

    List<String> iconList = [
      "assets/icon/github.png",
      "assets/icon/youtube.png",
      "assets/icon/spotify.png",
    ];

    int serviceIndex = Random().nextInt(colorList.length);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              color: colorList[serviceIndex],
              border: Border.all(
                color: colorList[serviceIndex],
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),
            child: InkWell(
              onTap: () {
                debugPrint('Card tapped.');
              },
              child: Column(
                children: [
                  SizedBox(
                    width: constraints.maxWidth,
                  ),
                  Padding(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 5, vertical: 30),
                    child: Column(children: [
                      Image.asset(iconList[serviceIndex],
                          width: 75, height: 75),
                      AutoSizeText(
                        nameList[serviceIndex],
                        style: const TextStyle(
                            fontSize: 25,
                            fontWeight: FontWeight.bold,
                            color: Colors.white),
                        maxLines: 1,
                      )
                    ]),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class WorkflowCard extends StatelessWidget {
  const WorkflowCard({super.key});

  @override
  Widget build(BuildContext context) {
    List<Color> colorList = [
      const Color.fromARGB(255, 46, 41, 78),
      const Color.fromARGB(255, 44, 110, 252),
      const Color.fromARGB(255, 228, 64, 95),
      const Color.fromARGB(255, 29, 185, 84),
      const Color.fromARGB(255, 205, 65, 43),
    ];

    int colorIndex = Random().nextInt(colorList.length);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: LayoutBuilder(
        builder: (context, constraints) {
          return Container(
            decoration: BoxDecoration(
              border: Border.all(
                color: colorList[colorIndex],
                width: 2.0,
              ),
              borderRadius: BorderRadius.circular(4.0),
            ),
            child: InkWell(
              onTap: () {
                debugPrint('Card tapped.');
              },
              child: Column(
                children: [
                  SizedBox(
                    width: constraints.maxWidth,
                    height: 20,
                    child: DecoratedBox(
                      decoration: BoxDecoration(
                        color: colorList[colorIndex],
                      ),
                    ),
                  ),
                  Padding(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 5, vertical: 10),
                    child: SizedBox(
                      width: constraints.maxWidth,
                      child: const Text(
                        "A card that can be tapped\nAnd how fit with his content",
                        style: TextStyle(
                          fontSize: 25,
                          fontWeight: FontWeight.bold,
                          color: Colors.black,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}

class MySearchBar extends StatelessWidget {
  const MySearchBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 15, vertical: 10),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(5),
        ),
        child: TextField(
          decoration: InputDecoration(
            prefixIcon: Padding(
              padding: EdgeInsets.only(left: 5),
              child: Icon(
                Icons.search,
                size: 35,
                color: Colors.black,
              ),
            ),
            hintStyle: TextStyle(
              fontSize: 20,
              color: const Color.fromARGB(255, 156, 156, 156),
              fontWeight: FontWeight.bold,
            ),
            hintText: 'Search',
            contentPadding: EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(10),
              borderSide: BorderSide(
                color: const Color.fromARGB(255, 238, 238, 238),
                width: 5,
              ),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(10),
              borderSide: BorderSide(
                color: Colors.black,
                width: 4,
              ),
            ),
          ),
        ),
      ),
    );
  }
}
