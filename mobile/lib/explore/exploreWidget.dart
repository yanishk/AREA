import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "../custom_widget.dart";
import "../globalData.dart";
import 'exploreDisplayActionsOrReactions.dart';

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

class MySquareServicesCard extends StatelessWidget {
  const MySquareServicesCard({
    super.key,
    required this.name,
    required this.color,
    this.imageUrl = "",
    this.index = 0,
  });

  final String name;
  final Color color;
  final String imageUrl;
  final int index;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
      child: Container(
        width: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        height: (MediaQuery.sizeOf(Globaldata.myContext).width / 2) - 25,
        decoration: BoxDecoration(
          color: color,
          border: Border.all(
            color: color,
            width: 2.0,
          ),
          borderRadius: BorderRadius.circular(6.0),
        ),
        child: InkWell(
          onTap: () {
            Navigator.push(Globaldata.myContext, MaterialPageRoute(builder: (context) => DisplayActionOrReactionPage(serviceIndex: index,)));
          },
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SafeWebImage(
                url: "${Globaldata.domainName}$imageUrl",
                width: 50,
                height: 50,
              ),
              AutoSizeText(
                name,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white
                ),
                maxLines: 1,
              )
            ],
          ),
        ),
      ),
    );
  }
}
